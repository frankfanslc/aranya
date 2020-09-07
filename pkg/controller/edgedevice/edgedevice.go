/*
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package edgedevice

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"

	"arhat.dev/pkg/log"
	"arhat.dev/pkg/queue"
	"arhat.dev/pkg/reconcile"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"

	"arhat.dev/aranya-proto/gopb"
	aranyaapi "arhat.dev/aranya/pkg/apis/aranya/v1alpha1"
	"arhat.dev/aranya/pkg/conf"
	"arhat.dev/aranya/pkg/constant"
	"arhat.dev/aranya/pkg/virtualnode"
	"arhat.dev/aranya/pkg/virtualnode/connectivity"
	"arhat.dev/aranya/pkg/virtualnode/device"
	"arhat.dev/aranya/pkg/virtualnode/metrics"
	"arhat.dev/aranya/pkg/virtualnode/pod"
	"arhat.dev/aranya/pkg/virtualnode/storage"
)

func (c *Controller) onEdgeDeviceAdded(obj interface{}) *reconcile.Result {
	var (
		err        error
		edgeDevice = obj.(*aranyaapi.EdgeDevice)
		name       = edgeDevice.Name
		logger     = c.Log.WithFields(log.String("name", name))
	)

	if edgeDevice.Status.HostNode != c.hostNodeName {
		ed := edgeDevice.DeepCopy()
		ed.Status.HostNode = c.hostNodeName

		logger.V("updating edge device status")
		_, err = c.edgeDeviceClient.UpdateStatus(c.Context(), ed, metav1.UpdateOptions{})
		if err != nil {
			logger.I("failed to update edge device status", log.Error(err))
			return &reconcile.Result{Err: err}
		}
	}

	err = c.requestNodeClusterRoleEnsure()
	if err != nil {
		logger.I("failed to request node cluster role ensure", log.Error(err))
		return &reconcile.Result{Err: err}
	}

	logger.V("scheduling edge device instantiation work")
	c.vnRec.Update(name, nil, edgeDevice.DeepCopy())
	err = c.vnRec.Schedule(queue.Job{Action: queue.ActionUpdate, Key: name}, 0)
	if err != nil {
		logger.I("failed to schedule edge device instantiation work", log.NamedError("reason", err))
		// trigger reschedule to make sure it will be created
		return &reconcile.Result{Err: err}
	}

	return nil
}

func (c *Controller) onEdgeDeviceUpdated(oldObj, newObj interface{}) *reconcile.Result {
	var (
		err           error
		oldEdgeDevice = oldObj.(*aranyaapi.EdgeDevice)
		newEdgeDevice = newObj.(*aranyaapi.EdgeDevice)
		name          = oldEdgeDevice.Name
		logger        = c.Log.WithFields(log.String("name", name))
	)

	if !c.checkEdgeDeviceUpdated(oldEdgeDevice, newEdgeDevice) {
		// no need to update virtualnode and related resources
		return nil
	}

	logger.V("scheduling edge device update work")
	c.vnRec.Update(name, oldEdgeDevice.DeepCopy(), newEdgeDevice.DeepCopy())
	err = c.vnRec.Schedule(queue.Job{Action: queue.ActionUpdate, Key: name}, 0)
	if err != nil {
		logger.I("failed to schedule edge device update work", log.NamedError("reason", err))
		return &reconcile.Result{Err: err}
	}

	return nil
}

func (c *Controller) onEdgeDeviceDeleting(obj interface{}) *reconcile.Result {
	var (
		edgeDevice = obj.(*aranyaapi.EdgeDevice)
		name       = edgeDevice.Name
		logger     = c.Log.WithFields(log.String("name", name))
	)

	logger.D("scheduling edge device deletion work", log.String("name", name))

	c.vnRec.Update(name, nil, edgeDevice.DeepCopy())
	err := c.vnRec.Schedule(queue.Job{Action: queue.ActionDelete, Key: name}, 0)
	if err != nil {
		logger.I("failed to schedule edge device deletion work", log.NamedError("reason", err))
		return &reconcile.Result{Err: err}
	}

	return nil
}

func (c *Controller) onEdgeDeviceDeleted(obj interface{}) *reconcile.Result {
	var (
		edgeDevice = obj.(*aranyaapi.EdgeDevice)
		name       = edgeDevice.Name
		logger     = c.Log.WithFields(log.String("name", name))
	)

	logger.V("requesting node cluster role update")
	err := c.requestNodeClusterRoleEnsure()
	if err != nil {
		logger.I("failed to request node cluster role update", log.Error(err))
		return &reconcile.Result{Err: err}
	}

	if _, ok := c.virtualNodes.Get(name); !ok {
		return nil
	}

	logger.V("scheduling edge device deletion work")
	c.vnRec.Update(name, nil, edgeDevice.DeepCopy())
	err = c.vnRec.Schedule(queue.Job{Action: queue.ActionDelete, Key: name}, 0)
	if err != nil {
		logger.I("failed to schedule edge device deletion work")
		// do not trigger reschedule since we can rely on job queue's deduplication behavior
		return nil
	}

	return nil
}

// instantiationEdgeDevice create resources in sequence
//   - Ensure kubelet certification valid (secret object contains valid certificate for this node)
//   - Create connectivity manager
//   - Create virtualnode
//   - Ensure NodeVerbs object
//	 - Start virtualnode
func (c *Controller) instantiationEdgeDevice(name string) (err error) {
	var (
		logger = c.Log.WithFields(log.String("name", name))
		opts   = &virtualnode.CreationOptions{
			NodeName: name,
			Hostname: c.hostname,
			HostIP:   c.hostIP,

			SSHPrivateKey:      make([]byte, len(c.sshPrivateKey)),
			EventBroadcaster:   record.NewBroadcaster(),
			VirtualnodeManager: c.virtualNodes,

			ScheduleNodeSync: func() error {
				err2 := c.nodeStatusRec.Schedule(queue.Job{Action: queue.ActionUpdate, Key: name}, 0)
				if err2 != nil && !errors.Is(err2, queue.ErrJobDuplicated) {
					return err2
				}

				return nil
			},

			PodOptions: nil,
		}
	)
	_ = copy(opts.SSHPrivateKey, c.sshPrivateKey)

	ed, ok := c.getEdgeDeviceObject(name)
	if !ok {
		logger.I("failed to get edge device object for creation")
		return nil
	}

	logger.D("ensuring kubelet cert for edge device node")
	kubeletCert, err := c.ensureKubeletServerCert(c.hostNodeName, name, ed.Spec.Node.Cert, c.hostNodeAddresses)
	if err != nil {
		logger.I("failed to get kubelet cert for edge device node", log.Error(err))
		return err
	}

	logger.D("creating listener for kubelet server")
	httpListener, err := net.Listen("tcp", ":0")
	if err != nil {
		logger.I("failed to create listener for kubelet server", log.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			_ = httpListener.Close()
		}
	}()

	// create a new kube client just in case of rate limiting affected by other client
	opts.KubeClient, _, err = c.vnConfig.KubeClient.NewKubeClient(c.vnKubeconfig, true)
	if err != nil {
		logger.I("failed to create kube client", log.Error(err))
		return err
	}

	opts.KubeletServerListener = tls.NewListener(
		httpListener, &tls.Config{Certificates: []tls.Certificate{*kubeletCert}})
	defer func() {
		if err != nil {
			logger.V("closing kubelet server listener due to error")
			_ = opts.KubeletServerListener.Close()
		}
	}()

	// prepare all required options
	{
		vnConfig := c.vnConfig.OverrideWith(ed.Spec)

		logger.D("preparing connectivity options")
		opts.ConnectivityOptions, err = c.prepareConnectivityOptions(ed, vnConfig)
		if err != nil {
			logger.I("failed to prepare connectivity options", log.Error(err))
			return err
		}

		logger.D("preparing node options")
		opts.NodeOptions, err = c.prepareNodeOptions(ed, vnConfig)
		if err != nil {
			logger.I("failed to prepare node options", log.Error(err))
			return err
		}

		logger.D("preparing associated device options")
		opts.DeviceOptions, err = c.prepareDeviceOptions(ed, vnConfig)
		if err != nil {
			logger.I("failed to prepare associated device options", log.Error(err))
			return err
		}

		logger.D("preparing metrics options")
		opts.MetricsOptions, err = c.prepareMetricsOptions(ed, vnConfig)
		if err != nil {
			logger.I("failed to prepare metrics options", log.Error(err))
			return err
		}

		logger.D("preparing pod options")
		opts.PodOptions, err = c.preparePodOptions(ed, vnConfig, opts.EventBroadcaster)
		if err != nil {
			logger.I("failed to prepare pod options", log.Error(err))
			return err
		}

		logger.D("preparing storage options")
		opts.StorageOptions, err = c.prepareStorageOptions(ed, vnConfig, opts.EventBroadcaster)
		if err != nil {
			logger.I("failed to prepare storage options", log.Error(err))
			return err
		}
	}

	virtualNodeCtx, cancelVirtualNode := context.WithCancel(c.Context())
	logger.D("creating connectivity manager")
	opts.ConnectivityManager, err = connectivity.NewManager(virtualNodeCtx, name, opts.ConnectivityOptions)
	if err != nil {
		cancelVirtualNode()
		logger.I("failed to create connectivity manager", log.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			logger.I("closing connectivity manager due to error")
			opts.ConnectivityManager.Close()
		}
	}()

	logger.D("creating virtual node")
	vn, err := virtualnode.CreateVirtualNode(virtualNodeCtx, cancelVirtualNode, opts)
	if err != nil {
		logger.I("failed to create virtual node", log.Error(err))
		return err
	}

	defer func() {
		if err != nil {
			vn.ForceClose()
		}
	}()

	// MUST add virtual node before request node ensure
	err = c.virtualNodes.Add(vn)
	if err != nil {
		logger.I("failed to add virtual node, removing existing one", log.Error(err))
		c.virtualNodes.Delete(name)

		if err = c.virtualNodes.Add(vn); err != nil {
			logger.E("failed to add virtual node, unrecoverable", log.Error(err))
			return err
		}
	}

	defer func() {
		if err != nil {
			c.virtualNodes.Delete(name)
		}
	}()

	logger.D("requesting node ensure")
	err = c.requestNodeEnsure(name)
	if err != nil {
		logger.I("failed to request node ensure", log.Error(err))
		return err
	}

	if c.vnConfig.Node.Lease.Enabled {
		logger.D("requesting node lease ensure")
		err = c.requestNodeLeaseEnsure(name)
		if err != nil {
			logger.I("failed to request node lease ensure", log.Error(err))
			return err
		}
	}

	logger.D("requesting virtual pod ensure")
	err = c.requestVirtualPodEnsure(name)
	if err != nil {
		logger.I("failed to request virtual pod ensure", log.Error(err))
		return err
	}

	if vn.ConnectivityServerListener() != nil && c.svcReqRec != nil {
		logger.D("requesting connectivity service ensure")
		err = c.requestConnectivityServiceEnsure()
		if err != nil {
			logger.I("failed to request connectivity service ensure", log.Error(err))
			return err
		}
	}

	logger.D("starting virtual node")
	err = vn.Start()
	if err != nil {
		logger.I("failed to start virtual node", log.Error(err))
		return err
	}

	return nil
}

func (c *Controller) checkEdgeDeviceUpdated(old, new *aranyaapi.EdgeDevice) bool {
	old = old.DeepCopy()
	// ignore reverse proxy change
	old.Spec.ReverseProxy = new.Spec.ReverseProxy
	return !reflect.DeepEqual(old.Spec, new.Spec)
}

func (c *Controller) prepareConnectivityOptions(
	edgeDevice *aranyaapi.EdgeDevice,
	vnConfig *conf.VirtualnodeConfig,
) (config *connectivity.Options, err error) {
	config = &connectivity.Options{
		UnarySessionTimeout: vnConfig.Connectivity.Timers.UnarySessionTimeout,
		GRPCOpts:            nil,
		MQTTOpts:            nil,
		AMQPOpts:            nil,
		AzureIoTHubOpts:     nil,
		GCPIoTCoreOpts:      nil,
	}

	switch edgeDevice.Spec.Connectivity.Method {
	case aranyaapi.GCPIoTCore:
		opts, err2 := c.generateGCPPubSubOpts(edgeDevice.Spec.Connectivity.GCPIoTCore)
		if err2 != nil {
			return nil, fmt.Errorf("failed to generate gcp pub/sub opts: %w", err2)
		}

		config.GCPIoTCoreOpts = opts
	case aranyaapi.AzureIoTHub:
		opts, err2 := c.generateAzureIoTHubOpts(edgeDevice.Spec.Connectivity.AzureIoTHub)
		if err2 != nil {
			return nil, fmt.Errorf("failed to generate azure iot hub opts: %w", err2)
		}

		config.AzureIoTHubOpts = opts
	case aranyaapi.GRPC:
		var (
			grpcConfig     = edgeDevice.Spec.Connectivity.GRPC
			grpcSrvOptions []grpc.ServerOption
			tlsConfig      *tls.Config
		)

		tlsConfig, err = c.createTLSConfigFromSecret(grpcConfig.TLSSecretRef)
		if err != nil {
			c.Log.I("failed to get grpc server tls secret", log.Error(err))
			return nil, err
		}

		if tlsConfig != nil {
			grpcSrvOptions = append(grpcSrvOptions, grpc.Creds(credentials.NewTLS(tlsConfig)))
		}

		var grpcListener net.Listener
		grpcListener, err = net.Listen("tcp", ":0")
		if err != nil {
			return nil, err
		}

		defer func() {
			if err != nil {
				_ = grpcListener.Close()
			}
		}()

		config.GRPCOpts = &connectivity.GRPCOpts{
			Server:   grpc.NewServer(grpcSrvOptions...),
			Listener: grpcListener,
		}
	case aranyaapi.MQTT:
		var (
			mqttConfig = edgeDevice.Spec.Connectivity.MQTT
			tlsConfig  *tls.Config
			username   []byte
			password   []byte
		)

		tlsConfig, err = c.createTLSConfigFromSecret(mqttConfig.TLSSecretRef)
		if err != nil {
			c.Log.I("failed to get mqtt client tls secret", log.Error(err))
			return nil, err
		}

		if userPassRef := mqttConfig.UserPassRef; userPassRef != nil && userPassRef.Name != "" {
			c.Log.D("fetching mqtt user pass")
			username, password, err = c.getUserPassFromSecret(userPassRef.Name)
			if err != nil {
				c.Log.I("failed to get mqtt client user pass", log.Error(err))
				return nil, err
			}
		}

		config.MQTTOpts = &connectivity.MQTTOpts{
			TLSConfig: tlsConfig,
			Username:  username,
			Password:  password,
			Config:    mqttConfig,
		}
	case aranyaapi.AMQP:
		var (
			amqpConfig = edgeDevice.Spec.Connectivity.AMQP
			tlsConfig  *tls.Config
			username   []byte
			password   []byte
		)

		tlsConfig, err = c.createTLSConfigFromSecret(amqpConfig.TLSSecretRef)
		if err != nil {
			c.Log.I("failed to create amqp client tls config", log.Error(err))
			return nil, err
		}

		if userPassRef := amqpConfig.UserPassRef; userPassRef != nil && userPassRef.Name != "" {
			c.Log.D("fetching amqp username and password")
			username, password, err = c.getUserPassFromSecret(userPassRef.Name)
			if err != nil {
				c.Log.I("failed to get amqp username and password", log.Error(err))
				return nil, err
			}
		}

		config.AMQPOpts = &connectivity.AMQPOpts{
			TLSConfig: tlsConfig,
			Username:  username,
			Password:  password,
			Config:    amqpConfig,
		}
	default:
		err = fmt.Errorf("unsupported connectivity method: %s", edgeDevice.Spec.Connectivity.Method)
		return
	}

	return
}

// nolint:unparam
func (c *Controller) prepareNodeOptions(
	edgeDevice *aranyaapi.EdgeDevice,
	vnConfig *conf.VirtualnodeConfig,
) (_ *virtualnode.Options, err error) {
	opts := &virtualnode.Options{
		ForceSyncInterval:      vnConfig.Node.Timers.ForceSyncInterval,
		MirrorNodeSyncInterval: vnConfig.Node.Timers.MirrorSyncInterval,
	}

	return opts, nil
}

// nolint:unparam
func (c *Controller) preparePodOptions(
	edgeDevice *aranyaapi.EdgeDevice,
	vnConfig *conf.VirtualnodeConfig,
	evb record.EventBroadcaster,
) (_ *pod.Options, err error) {
	opts := &pod.Options{
		Config: &vnConfig.Pod,
		GetNode: func() *corev1.Node {
			node, ok := c.getNodeObject(edgeDevice.Name)
			if !ok {
				return nil
			}

			return node
		},
		GetPod: func(name string) *corev1.Pod {
			po, ok := c.getWatchPodObject(name)
			if !ok {
				return nil
			}

			return po
		},
		GetConfigMap: func(name string) *corev1.ConfigMap {
			cm, ok := c.getWatchConfigMapObject(name)
			if !ok {
				return nil
			}

			return cm
		},
		GetSecret: func(name string) *corev1.Secret {
			secret, ok := c.getWatchSecretObject(name)
			if !ok {
				return nil
			}

			return secret
		},
		ListServices: c.listWatchServiceObjects,
		EventRecorder: evb.NewRecorder(
			scheme.Scheme,
			corev1.EventSource{Component: fmt.Sprintf("aranya-pod-manager,%s", edgeDevice.Name)},
		),
	}

	return opts, nil
}

// nolint:unparam
func (c *Controller) prepareDeviceOptions(
	edgeDevice *aranyaapi.EdgeDevice,
	vnConfig *conf.VirtualnodeConfig,
) (_ *device.Options, err error) {
	var (
		devices = make(map[string]*gopb.Device)
		devs    = make(map[string]struct{})
	)

	for _, d := range edgeDevice.Spec.Devices {
		if d.Name == "" || d.Name == constant.VirtualContainerNameHost {
			err = fmt.Errorf("invalid device name %q", d.Name)
			return
		}

		if msgs := validation.IsDNS1123Label(d.Name); len(msgs) > 0 {
			err = fmt.Errorf("device name %q is not a valid dns label: %s", d.Name, strings.Join(msgs, ", "))
			return
		}

		if _, ok := devs[d.Name]; ok {
			err = fmt.Errorf("duplicate device name %q", d.Name)
			return
		}

		devs[d.Name] = struct{}{}

		var mode gopb.DeviceConnectivity_Mode
		switch d.Connectivity.Mode {
		case aranyaapi.DeviceConnectivityModeServer:
			mode = gopb.DEVICE_CONNECTIVITY_MODE_SERVER
		case aranyaapi.DeviceConnectivityModeClient:
			mode = gopb.DEVICE_CONNECTIVITY_MODE_CLIENT
		default:
			mode = gopb.DEVICE_CONNECTIVITY_MODE_CLIENT
		}

		var operations []*gopb.DeviceOperation
		for _, op := range d.Operations {
			operations = append(operations, &gopb.DeviceOperation{
				Id:              op.Name,
				TransportParams: op.TransportParams,
			})
		}

		var uploadConnectivity *gopb.DeviceConnectivity
		if uc := d.UploadConnectivity; uc != nil {
			uploadConnectivity = &gopb.DeviceConnectivity{
				Transport: uc.Transport,
				Mode:      gopb.DEVICE_CONNECTIVITY_MODE_CLIENT,
				Target:    uc.Target,
				Params:    uc.Params,
				//Tls:       ,
			}
		}

		var deviceMetrics []*gopb.DeviceMetrics
		for _, m := range d.Metrics {
			var uploadMethod gopb.DeviceMetrics_DeviceMetricsUploadMethod
			switch m.UploadMethod {
			case aranyaapi.UploadWithArhatConnectivity:
				uploadMethod = gopb.UPLOAD_WITH_ARHAT_CONNECTIVITY
			case aranyaapi.UploadWithNodeMetrics:
				uploadMethod = gopb.UPLOAD_WITH_NODE_METRICS
			case aranyaapi.UploadWithStandaloneClient:
				if uploadConnectivity == nil {
					err = fmt.Errorf("no upload connectivity configured for standalone upload method")
					return
				}
				uploadMethod = gopb.UPLOAD_WITH_STANDALONE_CLIENT
			}

			deviceMetrics = append(deviceMetrics, &gopb.DeviceMetrics{
				Id:              m.Name,
				TransportParams: m.TransportParams,
				UploadMethod:    uploadMethod,
				UploadParams:    m.UploadParams,
			})
		}

		devices[d.Name] = &gopb.Device{
			Connectivity: &gopb.DeviceConnectivity{
				Transport: d.Connectivity.Transport,
				Mode:      mode,
				Target:    d.Connectivity.Target,
				Params:    d.Connectivity.Params,
				//Tls:
			},
			Operations:         operations,
			Metrics:            deviceMetrics,
			UploadConnectivity: uploadConnectivity,
		}
	}

	return &device.Options{Devices: devices}, nil
}

// nolint:unparam
func (c *Controller) prepareMetricsOptions(
	edgeDevice *aranyaapi.EdgeDevice,
	vnConfig *conf.VirtualnodeConfig,
) (_ *metrics.Options, err error) {
	opts := &metrics.Options{
		NodeMetrics:      &vnConfig.Node.Metrics,
		ContainerMetrics: &vnConfig.Pod.Metrics,
	}

	return opts, nil
}

// nolint:unparam
func (c *Controller) prepareStorageOptions(
	edgeDevice *aranyaapi.EdgeDevice,
	vnConfig *conf.VirtualnodeConfig,
	evb record.EventBroadcaster,
) (_ *storage.Options, err error) {
	if !vnConfig.Node.Storage.Enabled {
		return nil, nil
	}

	opts := &storage.Options{
		StorageRootDir:         vnConfig.Node.Storage.RootDir,
		KubeletRegistrationDir: vnConfig.Node.Storage.KubeletRegistrationDir,
		KubeletPluginsDir:      vnConfig.Node.Storage.KubeletPluginsDir,
		CSIDriverLister:        c.csiDriverLister,
		EventRecorder: evb.NewRecorder(
			scheme.Scheme,
			corev1.EventSource{Component: fmt.Sprintf("aranya-storage-manager,%s", edgeDevice.Name)},
		),
	}

	return opts, nil
}