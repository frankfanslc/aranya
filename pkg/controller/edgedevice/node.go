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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"

	"arhat.dev/pkg/log"
	"arhat.dev/pkg/patchhelper"
	"arhat.dev/pkg/queue"
	"arhat.dev/pkg/reconcile"
	"arhat.dev/pkg/textquery"
	"github.com/itchyny/gojq"
	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilnode "k8s.io/kubernetes/pkg/util/node"

	aranyaapi "arhat.dev/aranya/pkg/apis/aranya/v1alpha1"
	"arhat.dev/aranya/pkg/constant"
	"arhat.dev/aranya/pkg/util/convert"
	"arhat.dev/aranya/pkg/virtualnode"
)

func (c *Controller) checkNodeQualified(node *corev1.Node, kubeletListener net.Listener) bool {
	port, err := getListenerPort(kubeletListener)
	if err != nil {
		return false
	}

	if len(node.Labels) == 0 || node.Labels[constant.LabelRole] != constant.LabelRoleValueNode {
		return false
	}

	hasRequiredTaint := false
	for _, t := range node.Spec.Taints {
		if t.Key == constant.TaintKeyNamespace &&
			t.Value == constant.WatchNS() &&
			t.Effect == corev1.TaintEffectNoSchedule {
			hasRequiredTaint = true
		}
	}

	return hasRequiredTaint && node.Status.DaemonEndpoints.KubeletEndpoint.Port == port
}

func (c *Controller) requestNodeEnsure(name string) error {
	if c.nodeReqRec == nil {
		return nil
	}

	c.nodeReqRec.Update(name, name, name)
	err := c.nodeReqRec.Schedule(queue.Job{Action: queue.ActionAdd, Key: name}, 0)
	if err != nil && !errors.Is(err, queue.ErrJobDuplicated) {
		return fmt.Errorf("failed to schedule node ensure: %w", err)
	}

	return nil
}

func (c *Controller) onNodeEnsureRequested(obj interface{}) *reconcile.Result {
	var (
		name   = obj.(string)
		logger = c.Log.WithFields(log.String("name", name))
	)

	vn, ok := c.virtualNodes.Get(name)
	if !ok {
		_, shouldEnsure := c.getEdgeDeviceObject(name)
		if !shouldEnsure {
			return nil
		}

		return &reconcile.Result{Err: fmt.Errorf("waiting for virtual node")}
	}

	logger.V("ensuring node")
	err := c.ensureNodeObject(name, vn.KubeletServerListener())
	if err != nil {
		logger.I("failed to ensure node")
		return &reconcile.Result{Err: err}
	}
	logger.V("ensured node")

	return nil
}

func (c *Controller) onNodeAdd(obj interface{}) *reconcile.Result {
	var (
		err    error
		node   = obj.(*corev1.Node).DeepCopy()
		name   = node.Name
		logger = c.Log.WithFields(log.String("name", name))
	)

	logger.D("evaluating node field hooks")
	err = c.evalNodeFiledHooksAndUpdateNode(logger, node)
	if err != nil {
		return &reconcile.Result{Err: err}
	}

	return &reconcile.Result{NextAction: queue.ActionUpdate}
}

func (c *Controller) onNodeUpdated(oldObj, newObj interface{}) *reconcile.Result {
	var (
		err     error
		oldNode = oldObj.(*corev1.Node).DeepCopy()
		newNode = newObj.(*corev1.Node).DeepCopy()
		name    = newNode.Name
		logger  = c.Log.WithFields(log.String("name", name))
	)

	vn, ignore := c.doNodeResourcePreCheck(name)
	if ignore {
		return nil
	}

	listener := vn.KubeletServerListener()
	if listener == nil {
		logger.E("no kubelet listener found, error in virtualnode creation")
		return nil
	}

	if c.checkNodeQualified(newNode, listener) {
		logger.V("matches virtual node, setting pod cidr(s)")
		vn.SetPodCIDRs(convert.GetPodCIDRWithIPVersion(newNode.Spec.PodCIDR, newNode.Spec.PodCIDRs))
	} else {
		logger.D("updating outdated node")
		err = c.ensureNodeObject(name, listener)
		if err != nil {
			logger.I("failed to ensure node object")
			return &reconcile.Result{Err: err}
		}

		// reschedule node update to check node annotations
		return &reconcile.Result{NextAction: queue.ActionUpdate}
	}

	if reflect.DeepEqual(oldNode.Annotations, newNode.Annotations) &&
		reflect.DeepEqual(oldNode.Labels, newNode.Labels) {
		// no changes to annotation and label, so no need to check node field hook
		logger.V("node is up to date, skipping field hooks")
		return nil
	}

	logger.D("evaluating node field hooks")
	err = c.evalNodeFiledHooksAndUpdateNode(logger, newNode)
	if err != nil {
		return &reconcile.Result{Err: err}
	}

	return nil
}

func (c *Controller) onNodeDeleting(obj interface{}) *reconcile.Result {
	// TODO: add finalizer
	err := c.nodeClient.Delete(c.Context(), obj.(*corev1.Node).Name, *deleteAtOnce)
	if err != nil && !kubeerrors.IsNotFound(err) {
		return &reconcile.Result{Err: err}
	}

	return &reconcile.Result{NextAction: queue.ActionCleanup}
}

func (c *Controller) onNodeDeleted(obj interface{}) *reconcile.Result {
	var (
		name   = obj.(*corev1.Node).Name
		logger = c.Log.WithFields(log.String("name", name))
	)

	// TODO: NodeVerbs could be removed from cache by changing the label, which wont trigger onNodeDeleting
	vn, ignore := c.doNodeResourcePreCheck(name)
	if ignore {
		return nil
	}

	listener := vn.KubeletServerListener()
	if listener == nil {
		logger.E("no kubelet listener found, error in virtualnode creation")
		return nil
	}

	// virtualnode exists, this node should present
	err := c.ensureNodeObject(name, listener)
	if err != nil {
		logger.I("failed to restore node deleted", log.Error(err))
		return &reconcile.Result{Err: err}
	}

	return nil
}

func (c *Controller) evalNodeFiledHooksAndUpdateNode(logger log.Interface, nodeObj *corev1.Node) error {
	device, ok := c.getEdgeDeviceObject(nodeObj.Name)
	if !ok {
		logger.I("no EdgeDevice for this node")
		return nil
	}

	node := nodeObj.DeepCopy()
	nodeJSONBytes, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshal node object: %w", err)
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(nodeJSONBytes, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal node object as map: %w", err)
	}

	nodeNeedUpdate := false
	for _, h := range device.Spec.Node.FieldHooks {
		var query *gojq.Query
		query, err = gojq.Parse(h.Query)
		if err != nil {
			logger.E("bad field query", log.String("query", h.Query), log.Error(err))
			continue
		}

		result, found, err2 := textquery.RunQuery(query, data, nil)
		if err2 != nil {
			logger.E("failed to run query over node object", log.String("query", h.Query), log.Error(err2))
			continue
		}

		if !found {
			// nothing found, which means query failed
			continue
		}

		// found something, proceed

		targetValue := h.Value
		switch {
		case targetValue != "":
		// user already set value, do nothing
		case targetValue == "" && h.ValueExpression != "":
			targetValue, err = textquery.JQ(h.ValueExpression, result)
			if err != nil {
				logger.E("failed to eval value expression", log.String("exp", h.ValueExpression), log.Error(err))
			}
		}

		fp := h.TargetFieldPath

		switch {
		case strings.HasPrefix(fp, "metadata.labels['") && strings.HasSuffix(fp, "']"):
			key := strings.TrimSuffix(strings.TrimPrefix(fp, "metadata.labels['"), "']")

			if node.Labels[key] != targetValue {
				node.Labels[key] = targetValue
				nodeNeedUpdate = true
			}
		case strings.HasPrefix(fp, "metadata.annotations['") && strings.HasSuffix(fp, "']"):
			key := strings.TrimSuffix(strings.TrimPrefix(fp, "metadata.annotations['"), "']")

			if node.Annotations[key] != targetValue {
				node.Annotations[key] = targetValue
				nodeNeedUpdate = true
			}
		default:
			logger.E("unsupported target field", log.String("targetField", fp))
			continue
		}
	}

	if nodeNeedUpdate {
		logger.D("updating node object metadata for field hooks")
		err = patchhelper.TwoWayMergePatch(nodeObj.DeepCopy(), node, &corev1.Node{},
			func(patchData []byte) error {
				_, err2 := c.nodeClient.Patch(
					c.Context(), node.Name, types.StrategicMergePatchType, patchData, metav1.PatchOptions{})
				if err2 != nil {
					return err2
				}
				return nil
			},
		)

		if err != nil {
			logger.E("failed to update node object", log.Error(err))
			return err
		}
	}

	return nil
}

func (c *Controller) ensureNodeObject(name string, kubeletListener net.Listener) error {
	var create bool
	edgeDeviceObj, ok := c.getEdgeDeviceObject(name)
	if !ok {
		// EdgeDevice is the controller of node, and once edge device get deleted,
		// node will be deleted as well (by garbage-collector in kube-controller-manager)
		// user should use `kubectl delete --cascade=false` when deleting edge device
		// if they want to preserve node object
		return nil
	}

	kubeletPort, err := getListenerPort(kubeletListener)
	if err != nil {
		return err
	}

	node := c.newNodeForEdgeDevice(edgeDeviceObj, kubeletPort)

	oldNode, err := c.nodeClient.Get(c.Context(), name, metav1.GetOptions{})
	if err == nil {
		c.Log.I("found old node, patch to update")
		err = func() error {
			clone := oldNode.DeepCopy()
			if clone.Annotations == nil {
				clone.Annotations = node.Annotations
			} else {
				// keep old annotations
				for k, v := range node.Annotations {
					clone.Annotations[k] = v
				}
			}

			if clone.Labels == nil {
				clone.Labels = node.Labels
			} else {
				// keep old labels
				for k, v := range node.Labels {
					if k == constant.LabelKubeRole {
						// keep label kubernetes.io/role
						continue
					}

					clone.Labels[k] = v
				}
			}

			var appliedTaints []corev1.Taint
			if len(clone.Spec.Taints) == 0 {
				appliedTaints = node.Spec.Taints
			} else {
				for i, t := range node.Spec.Taints {
					if t.Key == constant.TaintKeyNamespace {
						// ensure default taint
						appliedTaints = append(appliedTaints, node.Spec.Taints[i])
						continue
					}

					found := false
					for j, oldT := range clone.Spec.Taints {
						if t.Key == oldT.Key {
							found = true
							// do not override user settings for node
							appliedTaints = append(appliedTaints, clone.Spec.Taints[j])
							break
						}
					}

					if !found {
						appliedTaints = append(appliedTaints, node.Spec.Taints[i])
					}
				}
			}

			clone.ClusterName = node.ClusterName
			clone.Spec.Taints = appliedTaints
			clone.Spec.ProviderID = node.Spec.ProviderID
			clone.OwnerReferences = node.OwnerReferences

			err = patchhelper.TwoWayMergePatch(oldNode, clone, &corev1.Node{}, func(patchData []byte) error {
				oldNode, err = c.nodeClient.Patch(
					c.Context(), name, types.StrategicMergePatchType, patchData, metav1.PatchOptions{})
				return err
			})
			if err != nil {
				return fmt.Errorf("failed to patch node spec: %w", err)
			}

			oldNode, err = c.nodeClient.Get(c.Context(), name, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get patched node: %w", err)
			}

			clone = oldNode.DeepCopy()
			clone.Status.Addresses = node.Status.Addresses
			clone.Status.DaemonEndpoints = node.Status.DaemonEndpoints

			_, _, err = utilnode.PatchNodeStatus(c.kubeClient.CoreV1(), types.NodeName(name), oldNode, clone)
			if err != nil {
				return fmt.Errorf("failed to patch node status: %w", err)
			}

			if len(node.Spec.PodCIDRs) > 0 {
				for i := 0; i < 3; i++ {
					err = utilnode.PatchNodeCIDRs(c.kubeClient, types.NodeName(name), node.Spec.PodCIDRs)
					if err == nil {
						return nil
					}
				}

				return fmt.Errorf("failed to update pod cidr: %w", err)
			}

			return nil
		}()

		if err != nil {
			c.Log.I("failed to patch update node object", log.Error(err))

			if c.vnConfig.Node.RecreateIfPatchFailed {
				c.Log.I("deleting node after patch failed to recreate")
				err = c.nodeClient.Delete(c.Context(), name, *deleteAtOnce)
				if err != nil && !kubeerrors.IsNotFound(err) {
					return fmt.Errorf("failed to delete old node: %w", err)
				}

				create = true
			}
		}
	} else {
		if !kubeerrors.IsNotFound(err) {
			return fmt.Errorf("failed to get old node: %w", err)
		}

		// old node not found
		create = true
	}

	if create {
		node, err = c.nodeClient.Create(c.Context(), node, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

// nolint:lll
func (c *Controller) newNodeForEdgeDevice(edgeDevice *aranyaapi.EdgeDevice, kubeletPort int32) *corev1.Node {
	labels := map[string]string{
		// this label can be overridden
		constant.LabelKubeRole: constant.WatchNS(),
	}

	for k, v := range edgeDevice.Spec.Node.Labels {
		labels[k] = v
	}

	// ensure following labels not overridden by user, they are special labels
	labels[constant.LabelName] = edgeDevice.Name
	labels[constant.LabelRole] = constant.LabelRoleValueNode
	labels[constant.LabelNamespace] = constant.WatchNS()
	labels[corev1.LabelHostname] = c.hostname

	// ensure taints
	taints := []corev1.Taint{{
		Key:    constant.TaintKeyNamespace,
		Value:  edgeDevice.Namespace,
		Effect: corev1.TaintEffectNoSchedule,
	}}
	for _, t := range edgeDevice.Spec.Node.Taints {
		if t.Key == constant.TaintKeyNamespace {
			// do not override default taints
			continue
		}

		taints = append(taints, corev1.Taint{
			Key:    t.Key,
			Value:  t.Value,
			Effect: t.Effect,
		})
	}

	podCIDRs := []string{edgeDevice.Spec.Pod.IPv4CIDR, edgeDevice.Spec.Pod.IPv6CIDR}
	switch {
	case podCIDRs[0] == "" && podCIDRs[1] == "":
		// no cidr
		podCIDRs = nil
	case podCIDRs[0] == "" && podCIDRs[1] != "":
		// has ipv6 cidr
		podCIDRs = podCIDRs[1:]
	case podCIDRs[0] != "" && podCIDRs[1] == "":
		// has ipv4 cidr
		podCIDRs = podCIDRs[:1]
	}
	var firstPodCIDR string
	if len(podCIDRs) != 0 {
		firstPodCIDR = podCIDRs[0]
	}

	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:        edgeDevice.Name,
			Annotations: edgeDevice.Spec.Node.Annotations,
			Labels:      labels,
			ClusterName: edgeDevice.ClusterName,
			// TODO: currently kube-controller-manager (v1.18.5) claims the owner doesn't exist
			//       even if we are sure the EdgeDevice exists and the uid is correct
			//OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(edgeDevice, controllerKind)},
		},
		Spec: corev1.NodeSpec{
			Taints:   taints,
			PodCIDR:  firstPodCIDR,
			PodCIDRs: podCIDRs,
			// providerID is the cloud resource id with vendor name prefix, it is usually used by cloud controller
			// manager (CCM) to determine whether the related resource exists
			// format: <cloud-vendor-name>://<resource-id>
			//
			// cloud managed computing resources usually derive its resource id via some cloud internal service
			// 		e.g. digital-ocean droplets can retrieve its resource id by requesting http://169.254.169.254/metadata/v1/id
			//			 (ref: https://github.com/digitalocean/digitalocean-cloud-controller-manager/blob/master/docs/getting-started.md#--provider-iddigitaloceandroplet-id)
			//
			// if we do not provide it at all, chances are that, the CCM will determine it automatically and delete
			// this node object automatically if related resource not found
			//
			// if we do not provide it with a prefix unknown to the CCM (e.g. `gce://`), it will be translated
			// into cloud provider specific resource id to look up such resource
			//
			// so here we just set it to `aranya://<edgedevice-namespace>/<edgedevice-name>`
			//
			// working CCM environment:
			//	- digital-ocean
			ProviderID: c.getNodeProviderID(edgeDevice.Namespace, edgeDevice.Name),
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				// in kubernetes v1.18.5, there is an assumption for Node Resource object in kube-controller-manager
				// that all node have condition NodeReady present, or will cause nil pointer panic in the controller
				// and this node will get deleted
				{Type: corev1.NodeReady, Status: corev1.ConditionUnknown},
				{Type: corev1.NodeMemoryPressure, Status: corev1.ConditionUnknown},
				{Type: corev1.NodeDiskPressure, Status: corev1.ConditionUnknown},
				{Type: corev1.NodePIDPressure, Status: corev1.ConditionUnknown},
				{Type: corev1.NodeNetworkUnavailable, Status: corev1.ConditionUnknown},
			},
			Addresses:       c.hostNodeAddresses,
			DaemonEndpoints: corev1.NodeDaemonEndpoints{KubeletEndpoint: corev1.DaemonEndpoint{Port: kubeletPort}},
		},
	}
}

func (c *Controller) doNodeResourcePreCheck(name string) (_ *virtualnode.VirtualNode, ignore bool) {
	// check if virtual node exist (relatively cheap operation)
	vn, ok := c.virtualNodes.Get(name)
	if !ok {
		// do not manage pending node
		return nil, true
	}

	_, ok = c.getEdgeDeviceObject(name)
	if !ok {
		// do not manage unknown node
		return nil, true
	}

	return vn, false
}
