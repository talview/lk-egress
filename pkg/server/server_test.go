// Copyright 2026 LiveKit, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func waitReturns(s *Server, timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		s.waitForActiveRequests()
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// A draining server has shutdown broken, which makes IsIdle false forever.
// The drain wait must still return once there are no active requests.
func TestWaitForActiveRequestsReturnsWhenShutdownBroken(t *testing.T) {
	s := &Server{}
	s.shutdown.Break()
	require.False(t, s.IsIdle())

	require.True(t, waitReturns(s, 3*time.Second))
}

func TestWaitForActiveRequestsBlocksUntilLastRequestEnds(t *testing.T) {
	s := &Server{}
	s.shutdown.Break()
	s.activeRequests.Inc()

	require.False(t, waitReturns(s, 1500*time.Millisecond))

	s.activeRequests.Dec()
	require.True(t, waitReturns(s, 3*time.Second))
}
