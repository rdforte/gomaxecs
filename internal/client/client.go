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

// Package client provides an HTTP client.
package client

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/rdforte/gomaxecs/internal/config"
)

// New returns a new Client.
func New(cfg config.Client) *Client {
	return &Client{
		client: &http.Client{
			Timeout: cfg.HTTPTimeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: cfg.DialTimeout,
				}).DialContext,
				MaxIdleConns:          cfg.MaxIdleConns,
				MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
				DisableKeepAlives:     cfg.DisableKeepAlives,
				IdleConnTimeout:       cfg.IdleConnTimeout,
				TLSHandshakeTimeout:   cfg.TLSHandshakeTimeout,
				ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
			},
		},
	}
}

// Client is an HTTP client.
type Client struct {
	client *http.Client
}

// Get performs an HTTP GET request.
func (c *Client) Get(ctx context.Context, url string) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP GET request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{res.StatusCode, body}, nil
}

type Response struct {
	StatusCode int
	Body       []byte
}
