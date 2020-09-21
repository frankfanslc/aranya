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
	"fmt"
	"io/ioutil"
	"sync"
	"testing"

	"arhat.dev/aranya-proto/aranyagopb"
	"github.com/stretchr/testify/assert"
)

func BenchmarkSession_deliverMsg(b *testing.B) {
	for _, size := range []int{1, 8, 256, 512, 1024, 64 * 1024, 256 * 1024, 1024 * 1024} {
		b.Run(fmt.Sprintf("normal msg payload size %d", size), func(b *testing.B) {
			s := newSession(1)
			data := make([]*aranyagopb.Msg, b.N)
			payload := make([]byte, size)
			for i := 0; i < b.N; i++ {
				data[i] = aranyagopb.NewMsg(aranyagopb.EMPTY, 1, uint64(i), i == b.N-1, payload)
			}

			msgCh := s.msgCh
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.deliverMsg(data[i])
			}
			<-msgCh
		})

		b.Run(fmt.Sprintf("stream msg payload size %d", size), func(b *testing.B) {
			s := newSession(1)
			data := make([]*aranyagopb.Msg, b.N)
			payload := make([]byte, size)
			for i := 0; i < b.N; i++ {
				data[i] = aranyagopb.NewMsg(aranyagopb.MSG_DATA_DEFAULT, 1, uint64(i), i == b.N-1, payload)
			}

			msgCh := s.msgCh
			wg := new(sync.WaitGroup)
			wg.Add(1)
			go func() {
				defer wg.Done()
				for m := range msgCh {
					d := m.(*Data)
					_, _ = ioutil.Discard.Write(d.Payload)
				}
			}()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.deliverMsg(data[i])
			}
			wg.Wait()
		})
	}
}

func TestSession_deliverMsg(t *testing.T) {
	// these are chunks for a NodeStatusMsg
	// sid: 1600661406479047299
	msgDataChunks := [][]byte{
		{0x0a, 0x0c, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x0a, 0x46, 0x0a, 0x06, 0x64, 0x61, 0x72, 0x77, 0x69, 0x6e},
		{0x0a, 0x0e, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x01, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x1a, 0x05, 0x61, 0x6d, 0x64, 0x36, 0x34, 0x22, 0x06, 0x31},
		{0x0a, 0x0e, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x02, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x38, 0x2e, 0x37, 0x2e, 0x30, 0x62, 0x24, 0x43, 0x36, 0x31},
		{0x0a, 0x0e, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x03, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x32, 0x35, 0x31, 0x32, 0x41, 0x2d, 0x30, 0x36, 0x30, 0x42},
		{0x0a, 0x0e, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x04, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x2d, 0x35, 0x36, 0x34, 0x41, 0x2d, 0x42, 0x44, 0x36, 0x30},
		{0x0a, 0x0e, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x05, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x2d, 0x35, 0x39, 0x45, 0x35, 0x46, 0x37, 0x43, 0x38, 0x37},
		{0x0a, 0x0e, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x06, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x36, 0x30, 0x30, 0xaa, 0x01, 0x06, 0x0a, 0x04, 0x6e, 0x6f},
		{0x0a, 0x0e, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x07, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x6e, 0x65, 0x12, 0x0f, 0x08, 0x04, 0x10, 0x80, 0x80, 0x80},
		{0x0a, 0x0e, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x08, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x80, 0x20, 0x18, 0x80, 0xc0, 0x99, 0xa2, 0xa6, 0x07, 0x1a},
		{0x0a, 0x0e, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x09, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x0c, 0x08, 0x01, 0x10, 0x01, 0x18, 0x01, 0x20, 0x01, 0x28},
		{0x0a, 0x0e, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x0a, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x01, 0x30, 0x01, 0x22, 0x1c, 0x10, 0x01, 0x18, 0x01, 0x2a},
		{0x0a, 0x0e, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x0b, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x16, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x61, 0x72, 0x68},
		{0x0a, 0x0e, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x0c, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x0a, 0x61, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x2f, 0x6f, 0x6e, 0x6c},
		{0x0a, 0x10, 0x08, 0x6f, 0x10, 0x83, 0xb5, 0xf1, 0xd6, 0x98, 0xa2, 0xac, 0x9b, 0x16, 0x18, 0x0d, 0x20, 0x01, 0x22, 0x05, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5a, 0x03, 0x69, 0x6e, 0x65},
	}

	seqSeries := [][]uint64{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
		{13, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0},
		{13, 1, 2, 3, 4, 5, 8, 7, 6, 9, 10, 11, 12, 0},
	}

	for i, serial := range seqSeries {
		t.Run(fmt.Sprintf("serial %d", i), func(t *testing.T) {
			s := newSession(1)
			msgCh := s.msgCh
			go func() {
				for j, seq := range serial {
					ck := msgDataChunks[seq]
					msg := new(aranyagopb.Msg)
					err := msg.Unmarshal(ck)
					if !assert.NoError(t, err) {
						assert.FailNow(t, fmt.Sprintf("failed to unamrshal msg data: %v", err))
					}

					delivered, complete := s.deliverMsg(msg)
					assert.True(t, delivered)

					if j == len(msgDataChunks)-1 {
						assert.True(t, complete)
					} else {
						assert.False(t, complete)
					}
				}
			}()

			msg := <-msgCh
			assert.IsType(t, &aranyagopb.Msg{}, msg)
			assert.Nil(t, s.msgCh)
		})
	}
}
