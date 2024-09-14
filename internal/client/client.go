package client

import (
	"fmt"
	"net"
	"net/http"

	"github.com/rdforte/gomax-ecs/internal/config"
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
func (c *Client) Get(url string) (*Response, error) {
	res, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return &Response{res}, nil
}

type Response struct {
	*http.Response
}
