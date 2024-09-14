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
