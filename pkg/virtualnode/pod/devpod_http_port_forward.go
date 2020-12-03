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
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"arhat.dev/aranya-proto/aranyagopb"
	"arhat.dev/aranya-proto/aranyagopb/runtimepb"
	"arhat.dev/pkg/iohelper"
	"arhat.dev/pkg/log"
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
	logger log.Interface,
	w http.ResponseWriter,
	r *http.Request,
	podUID string,
	opts *kubeletpf.V4Options,
	idleTimeout, streamCreationTimeout time.Duration,
	supportedProtocols []string,
) error {
	if wsstream.IsWebSocketRequest(r) {
		logger.V("serving websocket")
		return m.handleWebSocketPortForward(logger, w, r, podUID, idleTimeout, opts)
	}

	logger.V("serving spdy streams")
	return m.handleHTTPPortForward(logger, w, r, podUID, idleTimeout, streamCreationTimeout, supportedProtocols)
}

func (m *Manager) handleWebSocketPortForward(
	logger log.Interface,
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
			}()

			m.doPortForward(r.Context(), logger, wg, s)
		}()
	}

	return nil
}

func (m *Manager) handleHTTPPortForward(
	logger log.Interface,
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
	streamChan := make(chan httpstream.Stream, 16)

	up := spdy.NewResponseUpgrader()
	conn := up.UpgradeResponse(w, r, validateNewHTTPStream(streamChan))
	if conn == nil {
		return fmt.Errorf("failed to upgrade http stream connection")
	}
	conn.SetIdleTimeout(idleTimeout)

	wg := new(sync.WaitGroup)
	timeoutCheck := time.NewTimer(streamCreationTimeout / 5)

	streamPairs := make(map[string]*stream)
	for {
		select {
		case <-conn.CloseChan():
			logger.D("request closed, closing all streams")
			if !timeoutCheck.Stop() {
				<-timeoutCheck.C
			}

			for _, s := range streamPairs {
				s.close("")
			}

			_ = conn.Close()

			wg.Wait()

			logger.V("all streams exited")

			return nil
		case <-timeoutCheck.C:
			logger.V("detecting expired streams")
			timeoutCheck.Reset(streamCreationTimeout / 5)

			var toDelete []string
			now := time.Now()
			for reqID, p := range streamPairs {
				// stream pair created, check if its creation has timed out
				if now.After(p.creationFailAt) && !p.prepared() {
					logger.D("stream creation timeout", log.String("reqID", reqID))

					p.close("stream creation timeout")
					toDelete = append(toDelete, reqID)
				}
			}

			for _, id := range toDelete {
				delete(streamPairs, id)
			}
		case hs := <-streamChan:
			reqID := hs.Headers().Get(api.PortForwardRequestIDHeader)
			logger.V("new stream created", log.String("reqID", reqID))

			s, hasStreamPair := streamPairs[reqID]
			if !hasStreamPair {
				s = &stream{
					reqID: reqID,

					creationFailAt: time.Now().Add(streamCreationTimeout),
					podUID:         podUID,
					protocol:       "tcp",
				}
				streamPairs[reqID] = s
			}

			err = m.handleNewHTTPStream(r.Context(), logger.WithFields(log.String("reqID", s.reqID)), wg, s, hs)
			if err != nil {
				logger.D("invalid stream", log.Error(err))
				s.writeErr(err.Error())
			}
		}
	}
}

func (m *Manager) handleNewHTTPStream(
	ctx context.Context,
	logger log.Interface,
	wg *sync.WaitGroup,
	s *stream,
	hs httpstream.Stream,
) error {
	switch t := hs.Headers().Get(api.StreamType); t {
	case api.StreamTypeData:
		if s.data != nil {
			return fmt.Errorf("data stream already assigned")
		}

		logger.V("set data stream")

		portString := hs.Headers().Get(api.PortHeader)
		port, _ := strconv.ParseInt(portString, 10, 32)

		s.data = hs
		s.port = int32(port)
	case api.StreamTypeError:
		if s.error != nil {
			return fmt.Errorf("error stream already assigned")
		}

		logger.V("set error stream")

		s.error = hs
	default:
		logger.D("ignored stream", log.String("streamType", t))
		return nil
	}

	if !s.prepared() {
		return nil
	}

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()

			logger.D("stream finished")
			s.close("")
		}()

		m.doPortForward(ctx, logger, wg, s)
	}()

	return nil
}

func (m *Manager) doPortForward(ctx context.Context, logger log.Interface, wg *sync.WaitGroup, s *stream) {
	logger.D("forwarding to remote")

	pfCmd := &aranyagopb.PortForwardCmd{
		PodUid:  s.podUID,
		Network: s.protocol,
		Port:    s.port,
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
			logger.D("failed to marshal port-forward options", log.Error(err))
			s.close(fmt.Sprintf("prepare: %v", err))
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
		logger.D("failed to create session", log.Error(err))
		s.close(fmt.Sprintf("prepare: %v", err))
		return
	}

	// write routine is closed when read routine finished

	wg.Add(1)
	go func() {
		defer func() {
			s.close("")

			wg.Done()
		}()

		connectivity.HandleMessages(msgCh, func(msg *aranyagopb.Msg) (exit bool) {
			if msgErr := msg.GetError(); msgErr != nil {
				logger.I("error happened", log.Error(msgErr))
				s.close(fmt.Sprintf("error happened in remote: %s", msgErr.Error()))
			} else {
				logger.I("unexpected non data msg", log.Int32("kind", int32(msg.Kind)), log.Uint64("sid", sid))
				s.close(fmt.Sprintf("unexpected non data msg %d", msg.Kind))
			}

			return true
		}, func(dataMsg *connectivity.Data) (exit bool) {
			_, err2 := s.data.Write(dataMsg.Payload)
			if err2 != nil {
				logger.I("failed to write data", log.Error(err2))
				s.close(fmt.Sprintf("write: %v", err2))
				return true
			}

			return false
		}, connectivity.LogUnknownMessage(m.Log))
	}()

	// read routine

	var seq uint64
	defer func() {
		_, _, _, err2 := m.ConnectivityManager.PostData(sid, aranyagopb.CMD_DATA_UPSTREAM, nextSeq(&seq), true, nil)
		if err2 != nil {
			logger.D("failed to post read close cmd", log.Error(err2))
		}

		// close session with best effort
		_, _, _ = m.ConnectivityManager.PostCmd(sid, aranyagopb.CMD_SESSION_CLOSE, &aranyagopb.SessionCloseCmd{
			Sid: sid,
		})
	}()

	r := iohelper.NewTimeoutReader(s.data)
	go r.FallbackReading(ctx.Done())

	bufSize := m.ConnectivityManager.MaxPayloadSize()
	if bufSize > constant.MaxBufSize {
		bufSize = constant.MaxBufSize
	}

	buf := make([]byte, bufSize)
	for r.WaitForData(ctx.Done()) {
		data, shouldCopy, err2 := r.Read(constant.PortForwardStreamReadTimeout, buf)
		if err2 != nil {
			if len(data) == 0 && err2 != iohelper.ErrDeadlineExceeded {
				logger.I("failed to read data", log.Error(err2))
				s.writeErr(fmt.Sprintf("read: %v", err2))
				return
			}
		}

		if shouldCopy {
			data = make([]byte, len(data))
			_ = copy(data, buf[:len(data)])
		}

		// do not check returned last seq since we have limited the buffer size
		_, _, _, err2 = m.ConnectivityManager.PostData(
			sid, aranyagopb.CMD_DATA_UPSTREAM, nextSeq(&seq), false, data,
		)

		if err2 != nil {
			// TODO: shall we redo data post until successful?
			logger.I("failed to post data", log.Error(err2))
			s.writeErr(fmt.Sprintf("post: %v", err2))
			return
		}
	}
}

type stream struct {
	reqID string

	creationFailAt time.Time

	podUID   string
	protocol string
	port     int32

	data  io.ReadWriteCloser
	error io.WriteCloser
}

func (s *stream) prepared() bool {
	// s.podUID can be empty (for host port-forward)
	// s.protocol can be empty (default to tcp)
	return s.data != nil && s.error != nil && s.port > 0
}

func (s *stream) close(reason string) {
	if reason != "" {
		s.writeErr(fmt.Sprintf("[%s]: %s", s.reqID, reason))
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
func validateNewHTTPStream(streams chan<- httpstream.Stream) func(httpstream.Stream, <-chan struct{}) error {
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
