package config_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rdforte/gomaxecs/internal/config"
)

func TestConfig_LoadConfiguration(t *testing.T) {
	metaURIEnv := "ECS_CONTAINER_METADATA_URI_V4"
	uri := "mock-ecs-metadata-uri/"
	t.Setenv(metaURIEnv, uri)

	cfg := config.New()

	wantURI := "mock-ecs-metadata-uri"
	wantCfg := config.Config{
		ContainerMetadataURI: wantURI,
		TaskMetadataURI:      fmt.Sprintf("%s/task", wantURI),
		Client: config.Client{
			HTTPTimeout:           time.Second * 5,
			DialTimeout:           time.Second,
			MaxIdleConns:          1,
			MaxIdleConnsPerHost:   1,
			DisableKeepAlives:     false,
			IdleConnTimeout:       time.Second,
			TLSHandshakeTimeout:   time.Second,
			ResponseHeaderTimeout: time.Second,
		},
	}

	assert.Equal(t, wantCfg, cfg)
}
