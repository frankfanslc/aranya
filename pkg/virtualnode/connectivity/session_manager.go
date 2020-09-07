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

package connectivity

import (
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"

	"arhat.dev/pkg/queue"

	"arhat.dev/aranya-proto/gopb"
)

type session struct {
	parent *SessionManager

	epoch uint64

	closed bool
	msgCh  chan interface{}

	seq  uint64
	seqQ *queue.SeqQueue
	mu   *sync.Mutex
}

func (s *session) nextSeq() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	seq := s.seq
	s.seq++
	return seq
}

func (s *session) deliverMsg(msg *gopb.Msg) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		// session closed, no message shall be delivered
		return false
	}

	var reallyCompleted bool
	if s.seqQ == nil || msg.Kind != gopb.MSG_DATA {
		select {
		case s.msgCh <- msg:
		}
		reallyCompleted = msg.Completed
	} else {
		dataMsg := new(gopb.Data)
		_ = dataMsg.Unmarshal(msg.Body)

		if msg.Completed {
			reallyCompleted = s.seqQ.SetMaxSeq(dataMsg.Seq)
		}

		var out []interface{}
		out, reallyCompleted = s.seqQ.Offer(dataMsg.Seq, dataMsg.Data)
		for _, da := range out {
			select {
			case s.msgCh <- da.([]byte):
			}
		}
	}

	if reallyCompleted {
		s.closed = true
		if s.msgCh != nil {
			close(s.msgCh)
			s.msgCh = nil
		}

		s.parent.tq.Remove(msg.SessionId)

		s.parent.mu.Lock()
		delete(s.parent.m, msg.SessionId)
		s.parent.mu.Unlock()
	}

	return true
}

func (s *session) close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	s.closed = true
	if s.msgCh != nil {
		close(s.msgCh)
		s.msgCh = nil
	}
}

func NewSessionManager() *SessionManager {
	now := uint64(time.Now().UTC().UnixNano())
	return &SessionManager{
		m:  make(map[uint64]*session),
		tq: queue.NewTimeoutQueue(),
		mu: new(sync.RWMutex),

		sid:   now,
		epoch: now,
	}
}

type SessionManager struct {
	m  map[uint64]*session
	tq *queue.TimeoutQueue

	sid   uint64
	epoch uint64

	mu *sync.RWMutex
}

// nextSid is only intended to be used when s.mu has been locked
func (m *SessionManager) nextSid() uint64 {
	used := true
	for used {
		m.sid++
		_, used = m.m[m.sid]
	}
	return m.sid
}

func (m *SessionManager) Start(stopCh <-chan struct{}) error {
	m.tq.Start(stopCh)

	go func() {
		timeoutCh := m.tq.TakeCh()

		for {
			select {
			case session, more := <-timeoutCh:
				if !more {
					return
				}

				if sid, ok := session.Key.(uint64); ok {
					m.Dispatch(gopb.NewTimeoutErrorMsg(sid))
					m.Delete(sid)
				}
			case <-stopCh:
				for len(timeoutCh) > 0 {
					<-timeoutCh
				}
				return
			}
		}
	}()

	return nil
}

// Add or reuse a session
func (m *SessionManager) Add(
	sid uint64,
	cmd proto.Marshaler,
	timeout time.Duration,
	expectSeqData bool,
) (realSid uint64, ch chan interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	realSid = sid
	// check if old sid valid
	if oldSession, ok := m.m[realSid]; ok {
		// reuse old channel if valid
		ch = oldSession.msgCh

		// nolint:gocritic
		switch c := cmd.(type) {
		case *gopb.PodOperationCmd:
			if opts, ok := c.Options.(*gopb.PodOperationCmd_InputOptions); ok {
				opts.InputOptions.Seq = oldSession.nextSeq()
			}
		}
	} else {
		// create new channel if invalid
		session := &session{
			parent: m,
			msgCh:  make(chan interface{}, 1),
			epoch:  m.epoch,
			mu:     new(sync.Mutex),
		}

		if expectSeqData {
			session.seqQ = queue.NewSeqQueue()
		}

		ch = session.msgCh
		realSid = m.nextSid()

		if timeout > 0 {
			_ = m.tq.OfferWithDelay(realSid, session, timeout)
		}
		m.m[realSid] = session
	}

	return realSid, ch
}

func (m *SessionManager) Dispatch(msg *gopb.Msg) bool {
	m.mu.RLock()
	session, ok := m.m[msg.SessionId]
	m.mu.RUnlock()

	// only deliver msg for this epoch or stream message
	if ok && (session.epoch == m.epoch || session.seqQ != nil) {
		return session.deliverMsg(msg)
	}

	m.Delete(msg.SessionId)
	return false
}

func (m *SessionManager) Delete(sid uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, ok := m.m[sid]; ok {
		session.close()
		m.tq.Remove(sid)
		delete(m.m, sid)
	}
}

func (m *SessionManager) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	allSid := make([]uint64, len(m.m))
	i := 0
	for sid, session := range m.m {
		if session.seqQ != nil {
			// do not close streams
			continue
		}

		session.close()
		m.tq.Remove(sid)
		allSid[i] = sid
		i++
	}

	for _, sid := range allSid {
		delete(m.m, sid)
	}

	now := uint64(time.Now().UTC().UnixNano())

	m.epoch = now
	// m.sid > now is the rare case, but we have to make sure sid is
	// always increasing in this session manager
	if m.sid < now {
		m.sid = now
	}
}

func (m *SessionManager) Remains() []uint64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.m) == 0 {
		return nil
	}

	result := make([]uint64, len(m.m))
	i := 0
	for sid := range m.m {
		result[i] = sid
		i++
	}
	return result
}

func (m *SessionManager) TimeoutRemains() []queue.TimeoutData {
	return m.tq.Remains()
}