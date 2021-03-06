// Copyright (c) 2015 Uber Technologies, Inc.
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tchannel_test

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	. "github.com/uber/tchannel/golang"

	"github.com/stretchr/testify/assert"
	"github.com/uber/tchannel/golang/raw"
	"github.com/uber/tchannel/golang/testutils"
)

var cn = "hello"

func TestWrapContextForTest(t *testing.T) {
	call := testutils.NewIncomingCall(cn)
	ctx, cancel := NewContext(time.Second)
	defer cancel()
	actual := WrapContextForTest(ctx, call)
	assert.Equal(t, call, CurrentCall(actual), "Incorrect call object returned.")
}

func TestShardKeyPropagates(t *testing.T) {
	WithVerifiedServer(t, nil, func(ch *Channel, hostPort string) {
		peerInfo := ch.PeerInfo()
		testutils.RegisterFunc(t, ch, "test", func(ctx context.Context, args *raw.Args) (*raw.Res, error) {
			return &raw.Res{
				Arg3: []byte(CurrentCall(ctx).ShardKey()),
			}, nil
		})

		ctx, cancel := NewContextBuilder(time.Second).Build()
		defer cancel()
		_, arg3, _, err := raw.Call(ctx, ch, peerInfo.HostPort, peerInfo.ServiceName, "test", nil, nil)
		assert.NoError(t, err, "Call failed")
		assert.Equal(t, arg3, []byte(""))

		ctx, cancel = NewContextBuilder(time.Second).
			SetShardKey("shard").Build()
		defer cancel()
		_, arg3, _, err = raw.Call(ctx, ch, peerInfo.HostPort, peerInfo.ServiceName, "test", nil, nil)
		assert.NoError(t, err, "Call failed")
		assert.Equal(t, string(arg3), "shard")
	})
}

func TestCurrentCallWithNilResult(t *testing.T) {
	ctx, cancel := NewContext(time.Second)
	defer cancel()
	call := CurrentCall(ctx)
	assert.Nil(t, call, "Should return nil.")
}
