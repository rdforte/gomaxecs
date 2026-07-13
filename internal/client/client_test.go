// Copyright 2004 Ryan Forte
//
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

package client_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/rdforte/gomaxecs/internal/client"
	"github.com/rdforte/gomaxecs/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubDoer implements client.HTTPDoer and returns a fixed *http.Response and
// error for every request, so we can exercise Client.Get without a live server.
type stubDoer struct {
	resp *http.Response
	err  error
}

func (s *stubDoer) Do(*http.Request) (*http.Response, error) {
	return s.resp, s.err
}

func newClientWithDoer(t *testing.T, doer client.HTTPDoer) *client.Client {
	t.Helper()
	cfg := config.New(config.WithLogger(func(string, ...any) {}))
	return client.NewWithClient(cfg.DebugLogf, doer)
}

func TestClient_Get_ReturnsErrorOnNilResponse(t *testing.T) {
	t.Parallel()

	// A custom client.HTTPDoer that returns (nil, nil) — the stdlib *http.Client
	// itself rejects this, so the guard in Get only fires for non-stdlib clients.
	c := newClientWithDoer(t, &stubDoer{resp: nil, err: nil})

	res, err := c.Get(context.Background(), "http://example.test/metadata")
	require.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "received nil response from HTTP client")
}

func TestClient_Get_ReturnsErrorWhenTransportFails(t *testing.T) {
	t.Parallel()

	wantErr := io.ErrUnexpectedEOF
	c := newClientWithDoer(t, &stubDoer{resp: nil, err: wantErr})

	res, err := c.Get(context.Background(), "http://example.test/metadata")
	require.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "failed to perform HTTP GET request")
}

func TestClient_Get_ReturnsResponseOnSuccess(t *testing.T) {
	t.Parallel()

	body := strings.NewReader("hello")
	c := newClientWithDoer(t, &stubDoer{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(body),
		},
	})

	res, err := c.Get(context.Background(), "http://example.test/metadata")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, []byte("hello"), res.Body)
}
