package config_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rdforte/gomaxecs/internal/config"
)

func Test_Config_LoadConfiguration(t *testing.T) {
	metaURIEnv := "ECS_CONTAINER_METADATA_URI_V4"
	containerID := "container-id"
	uri := strings.Join([]string{"mock-ecs-metadata-uri", "/", containerID}, "")
	t.Setenv(metaURIEnv, uri)

	cfg := config.New()

	wantCfg := config.Config{
		MetadataURI: uri,
		ConainerID:  containerID,
		Client: config.Client{
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

	assert.Equal(t, wantCfg, cfg)
}
