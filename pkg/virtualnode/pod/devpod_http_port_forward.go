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

package pod

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"arhat.dev/aranya-proto/aranyagopb"
	"arhat.dev/aranya-proto/aranyagopb/runtimepb"
	"arhat.dev/pkg/iohelper"
	"github.com/gogo/protobuf/proto"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/apimachinery/pkg/util/httpstream/spdy"
	"k8s.io/apiserver/pkg/util/wsstream"
	api "k8s.io/kubernetes/pkg/apis/core"
	kubeletpf "k8s.io/kubernetes/pkg/kubelet/cri/streaming/portforward"

	"arhat.dev/aranya/pkg/constant"
	"arhat.dev/aranya/pkg/virtualnode/connectivity"
)

const (
	v4BinaryWebsocketProtocol = "v4." + wsstream.ChannelWebSocketProtocol
	v4Base64WebsocketProtocol = "v4." + wsstream.Base64ChannelWebSocketProtocol
)

func (m *Manager) servePortForward(
	w http.ResponseWriter,
	r *http.Request,
	podUID string,
	opts *kubeletpf.V4Options,
	idleTimeout, streamCreationTimeout time.Duration,
	supportedProtocols []string,
) error {
	if wsstream.IsWebSocketRequest(r) {
		return m.handleWebSocketPortForward(w, r, podUID, idleTimeout, opts)
	}

	return m.handleHTTPPortForward(w, r, podUID, idleTimeout, streamCreationTimeout, supportedProtocols)
}

func (m *Manager) handleWebSocketPortForward(
	w http.ResponseWriter,
	r *http.Request,
	podUID string,
	idleTimeout time.Duration,
	opts *kubeletpf.V4Options,
) error {
	channels := make([]wsstream.ChannelType, 0, len(opts.Ports)*2)
	for i := 0; i < len(opts.Ports); i++ {
		channels = append(channels, wsstream.ReadWriteChannel, wsstream.WriteChannel)
	}
	conn := wsstream.NewConn(map[string]wsstream.ChannelProtocolConfig{
		"": {
			Binary:   true,
			Channels: channels,
		},
		v4BinaryWebsocketProtocol: {
			Binary:   true,
			Channels: channels,
		},
		v4Base64WebsocketProtocol: {
			Binary:   false,
			Channels: channels,
		},
	})
	conn.SetIdleTimeout(idleTimeout)

	_, streams, err := conn.Open(w, r)
	if err != nil {
		err = fmt.Errorf("failed to upgrade websocket connection: %w", err)
		return err
	}

	wg := new(sync.WaitGroup)
	defer func() {
		wg.Wait()
		_ = conn.Close()
	}()

	for i := range opts.Ports {
		s := &stream{
			podUID:   podUID,
			protocol: "tcp",
			port:     opts.Ports[1],

			data:  streams[i*2],
			error: streams[i*2+1],
		}

		portBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(portBytes, uint16(s.port))
		_, _ = s.data.Write(portBytes)
		_, _ = s.error.Write(portBytes)

		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()

				s.close("")
			}()

			m.doPortForward(s)
		}()
	}

	return nil
}

func (m *Manager) handleHTTPPortForward(
	w http.ResponseWriter,
	r *http.Request,
	podUID string,
	idleTimeout, streamCreationTimeout time.Duration,
	supportedProtocols []string,
) error {
	_, err := httpstream.Handshake(r, w, supportedProtocols)
	// negotiated protocol isn't currently used server side, but could be in the future
	if err != nil {
		// Handshake writes the error to the client
		return err
	}
	streamChan := make(chan httpstream.Stream, 1)

	up := spdy.NewResponseUpgrader()
	conn := up.UpgradeResponse(w, r, validateNewHTTPStream(streamChan))
	if conn == nil {
		return fmt.Errorf("failed to upgrade http stream connection")
	}
	conn.SetIdleTimeout(idleTimeout)

	wg := new(sync.WaitGroup)
	timeoutCheckTk := time.NewTicker(streamCreationTimeout / 5)

	streamPairs := make(map[string]*stream)
	for {
		select {
		case <-conn.CloseChan():
			// connection closed, close all streams
			timeoutCheckTk.Stop()

			for _, s := range streamPairs {
				s.close("")
			}

			_ = conn.Close()

			wg.Wait()

			return nil
		case <-timeoutCheckTk.C:
			var toDelete []string
			now := time.Now()
			for reqID, p := range streamPairs {
				// stream pair created, check if its creation has timed out
				if now.After(p.creationFailAt) && !p.prepared() {
					p.close("stream creation timeout")
					toDelete = append(toDelete, reqID)
				}
			}

			for _, id := range toDelete {
				delete(streamPairs, id)
			}
		case hs := <-streamChan:
			reqID := hs.Headers().Get(api.PortForwardRequestIDHeader)
			s, hasStreamPair := streamPairs[reqID]
			if !hasStreamPair {
				s = &stream{
					creationFailAt: time.Now().Add(streamCreationTimeout),
					podUID:         podUID,
					protocol:       "tcp",
				}
				streamPairs[reqID] = s
			}

			err = m.handleNewHTTPStream(wg, s, hs)
			if err != nil {
				s.writeErr(err.Error())
			}
		}
	}
}

func (m *Manager) handleNewHTTPStream(wg *sync.WaitGroup, s *stream, hs httpstream.Stream) error {
	switch hs.Headers().Get(api.StreamType) {
	case api.StreamTypeData:
		if s.data != nil {
			return fmt.Errorf("data stream already assigned")
		}

		portString := hs.Headers().Get(api.PortHeader)
		port, _ := strconv.ParseInt(portString, 10, 32)

		s.data = hs
		s.port = int32(port)
	case api.StreamTypeError:
		if s.error != nil {
			return fmt.Errorf("error stream already assigned")
		}

		s.error = hs
	default:
		// ignore other streams
		return nil
	}

	if !s.prepared() {
		return nil
	}

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()

			s.close("")
		}()

		m.doPortForward(s)
	}()

	return nil
}

func (m *Manager) doPortForward(s *stream) {
	pfCmd := &aranyagopb.PortForwardCmd{
		PodUid:   s.podUID,
		Port:     s.port,
		Protocol: s.protocol,
	}

	var (
		kind aranyagopb.CmdType
		cmd  proto.Marshaler
	)

	if s.podUID == "" {
		// is host port-forward
		kind = aranyagopb.CMD_PORT_FORWARD
		cmd = pfCmd
	} else {
		// is runtime pod port-forward
		data, err := pfCmd.Marshal()
		if err != nil {
			s.writeErr(fmt.Sprintf("failed to marshal port-forward options: %v", err))
			return
		}

		kind = aranyagopb.CMD_RUNTIME
		cmd = &runtimepb.Packet{
			Kind:    runtimepb.CMD_PORT_FORWARD,
			Payload: data,
		}
	}

	msgCh, sid, err := m.ConnectivityManager.PostCmd(0, kind, cmd)
	if err != nil {
		s.writeErr(fmt.Sprintf("failed to create port-forward session: %v", err))
		return
	}

	var seq uint64
	go func() {
		defer func() {
			_, _, _, err2 := m.ConnectivityManager.PostData(sid, aranyagopb.CMD_DATA_UPSTREAM, nextSeq(&seq), true, nil)
			if err2 != nil {
				s.writeErr(fmt.Sprintf("failed to post port-forward read close cmd: %v", err2))
			}
		}()

		r := iohelper.NewTimeoutReader(s.data)
		go r.FallbackReading(m.Context().Done())

		bufSize := m.ConnectivityManager.MaxPayloadSize()
		if bufSize > constant.MaxBufSize {
			bufSize = constant.MaxBufSize
		}
		buf := make([]byte, bufSize)
		for r.WaitForData(m.Context().Done()) {
			data, shouldCopy, err2 := r.Read(constant.DefaultPortForwardStreamReadTimeout, buf)
			if err2 != nil {
				if len(data) == 0 && err2 != iohelper.ErrDeadlineExceeded {
					break
				}
			}

			if shouldCopy {
				data = make([]byte, len(data))
				_ = copy(data, buf[:len(data)])
			}

			_, _, lastSeq, err2 := m.ConnectivityManager.PostData(
				sid, aranyagopb.CMD_DATA_UPSTREAM, nextSeq(&seq), false, data,
			)

			atomic.StoreUint64(&seq, lastSeq+1)
			if err2 != nil {
				s.writeErr(err2.Error())
			}
		}
	}()

	connectivity.HandleMessages(msgCh, func(msg *aranyagopb.Msg) (exit bool) {
		if msgErr := msg.GetError(); msgErr != nil {
			s.writeErr(msgErr.Error())
		} else {
			s.writeErr(fmt.Sprintf("unexpected non data msg in session %d", sid))
		}

		return false
	}, func(dataMsg *connectivity.Data) (exit bool) {
		_, err2 := s.data.Write(dataMsg.Payload)
		if err2 != nil && !errors.Is(err2, io.EOF) {
			s.writeErr(err2.Error())
			return true
		}

		return false
	}, connectivity.LogUnknownMessage(m.Log))
}

type stream struct {
	creationFailAt time.Time

	podUID   string
	protocol string
	port     int32

	data  io.ReadWriteCloser
	error io.WriteCloser
}

func (s *stream) prepared() bool {
	// s.podUID can be empty (for host port-forward)
	return s.data != nil && s.error != nil && s.port > 0
}

func (s *stream) close(reason string) {
	if reason != "" {
		s.writeErr(reason)
	}

	if s.data != nil {
		_ = s.data.Close()
	}

	if s.error != nil {
		_ = s.error.Close()
	}
}

func (s *stream) writeErr(err string) {
	if s.error != nil {
		_, _ = s.error.Write([]byte(err))
	}
}

// validateNewHTTPStream is the httpstream.NewStreamHandler for port
// forward streams. It checks each stream's port and stream type headers,
// rejecting any streams that with missing or invalid values. Each valid
// stream is sent to the streams channel.
func validateNewHTTPStream(streams chan httpstream.Stream) func(httpstream.Stream, <-chan struct{}) error {
	return func(stream httpstream.Stream, replySent <-chan struct{}) error {
		// make sure it has a valid port header
		portString := stream.Headers().Get(api.PortHeader)
		if len(portString) == 0 {
			return fmt.Errorf("%q header is required", api.PortHeader)
		}
		port, err := strconv.ParseUint(portString, 10, 16)
		if err != nil {
			return fmt.Errorf("unable to parse %q as a port: %v", portString, err)
		}
		if port < 1 {
			return fmt.Errorf("port %q must be > 0", portString)
		}

		// make sure it has a valid stream type header
		streamType := stream.Headers().Get(api.StreamType)
		if len(streamType) == 0 {
			return fmt.Errorf("%q header is required", api.StreamType)
		}
		if streamType != api.StreamTypeError && streamType != api.StreamTypeData {
			return fmt.Errorf("invalid stream type %q", streamType)
		}

		streams <- stream
		return nil
	}
}
