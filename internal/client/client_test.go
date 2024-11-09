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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rdforte/gomaxecs/internal/client"
	"github.com/rdforte/gomaxecs/internal/config"
)

func TestClient_Get_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	cfg := config.Client{}
	c := client.New(cfg)

	_, err := c.Get(context.Background(), ts.URL)
	assert.NoError(t, err)
}

func TestClient_Get_BuildRequestFailure(t *testing.T) {
	t.Parallel()

	cfg := config.Client{}
	c := client.New(cfg)

	_, err := c.Get(context.Background(), "://invalid-url")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create HTTP request")
}

func TestClient_Get_ClientFailure(t *testing.T) {
	t.Parallel()

	cfg := config.Client{}
	c := client.New(cfg)

	_, err := c.Get(context.Background(), "invalid-url")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to perform HTTP GET request")
}

func TestClient_Get_ResBodyFailure(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("partial-data"))
		assert.NoError(t, err)

		if hijacker, ok := w.(http.Hijacker); ok {
			conn, _, _ := hijacker.Hijack()
			conn.Close()
		}
	}))

	cfg := config.Client{}
	c := client.New(cfg)

	_, err := c.Get(context.Background(), ts.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read response body")
}
