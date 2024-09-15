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

// Package config provides the package configuration.
package config

import (
	"os"
	"strings"
	"time"
)

const metaURIEnv = "ECS_CONTAINER_METADATA_URI_V4"

func New() Config {
	uri := metadataURI()

	pathParts := strings.Split(uri, "/")
	containerID := pathParts[len(pathParts)-1]

	return Config{
		MetadataURI: uri,
		ConainerID:  containerID,
		Client: Client{
			HTTPTimeout:           time.Second * 5,
			DialTimeout:           time.Second,
			MaxIdleConns:          1,
			MaxIdleConnsPerHost:   1,
			DisableKeepAlives:     true,
			IdleConnTimeout:       time.Second,
			TLSHandshakeTimeout:   time.Second,
			ResponseHeaderTimeout: time.Second,
		},
	}
}

func metadataURI() string {
	return os.Getenv(metaURIEnv)
}

// Config represents the packagge configuration.
type Config struct {
	MetadataURI string
	ConainerID  string
	Client      Client
}

// Client represents the HTTP client configuration.
type Client struct {
	HTTPTimeout           time.Duration
	DialTimeout           time.Duration
	MaxIdleConns          int
	MaxIdleConnsPerHost   int
	DisableKeepAlives     bool
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ResponseHeaderTimeout time.Duration
}
