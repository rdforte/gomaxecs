package config_test

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rdforte/gomaxecs/internal/config"
)

func TestConfig_New_LoadConfiguration(t *testing.T) {
	metaURIEnv := "ECS_CONTAINER_METADATA_URI_V4"
	uri := "mock-ecs-metadata-uri/"
	t.Setenv(metaURIEnv, uri)

	cfg := config.New()

	wantURI := "mock-ecs-metadata-uri"
	wantCfg := config.Config{
		ContainerMetadataURI: wantURI,
		TaskMetadataURI:      wantURI + "/task",
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

func TestConfig_New_AppliesOptions(t *testing.T) {
	opt1 := mockOption{}
	opt2 := mockOption{}

	config.New(&opt1, &opt2)

	assert.True(t, opt1.isApplied)
	assert.True(t, opt2.isApplied)
}

func TestConfig_WithLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	cfg := config.New(config.WithLogger(logger.Printf))

	cfg.Log("test log: %s, %s", "arg1", "arg2")

	wantLog := "test log: arg1, arg2\n"
	assert.Equal(t, wantLog, buf.String())
}

type mockOption struct {
	isApplied bool
}

func (m *mockOption) Apply(_ *config.Config) {
	m.isApplied = true
}
