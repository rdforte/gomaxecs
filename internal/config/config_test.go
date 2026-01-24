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
	t.Parallel()

	opt1 := mockOption{}
	opt2 := mockOption{}

	config.New(opt1.Apply, opt2.Apply)

	assert.True(t, opt1.isApplied)
	assert.True(t, opt2.isApplied)
}

func TestConfig_WithLogger_LogsMessage(t *testing.T) {
	t.Parallel()

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	cfg := config.New(config.WithLogger(logger.Printf))

	cfg.Logf("test log: %s, %s", "arg1", "arg2")

	wantLog := "test log: arg1, arg2\n"
	assert.Equal(t, wantLog, buf.String())
}

func TestConfig_GetECSMetadataURI_RetrievesMetadataURIFromEnv(t *testing.T) {
	metaURIEnv := "ECS_CONTAINER_METADATA_URI_V4"
	uri := "mock-ecs-metadata-uri/"
	t.Setenv(metaURIEnv, uri)

	got := config.GetECSMetadataURI()

	want := "mock-ecs-metadata-uri"
	assert.Equal(t, want, got)
}

func TestConfig_GomaxecsDebugEnv_SetsDebugModeInConfig(t *testing.T) {
	metaURIEnv := "ECS_CONTAINER_METADATA_URI_V4"
	uri := "mock-ecs-metadata-uri/"
	t.Setenv(metaURIEnv, uri)

	got := config.GetECSMetadataURI()

	want := "mock-ecs-metadata-uri"
	assert.Equal(t, want, got)
}

func TestConfig_DebugLogf_LogsWhenGomaxecsDebugEnvEnabled(t *testing.T) {
	tableTest := []struct {
		env string
	}{
		{env: "true"},
		{env: "1"},
	}

	for _, tt := range tableTest {
		t.Run("GOMAXECS_DEBUG="+tt.env, func(t *testing.T) {
			t.Setenv("GOMAXECS_DEBUG", tt.env)

			buf := new(bytes.Buffer)
			logger := log.New(buf, "", 0)

			cfg := config.New(config.WithLogger(logger.Printf))

			// clear buffer so we only have fresh logs
			buf.Reset()

			cfg.DebugLogf("debug log: %s", "stub-message")

			wantLog := "gomaxecs debug: debug log: stub-message\n"
			assert.Equal(t, wantLog, buf.String())
		})
	}
}

func TestConfig_DebugLogf_DoesNotLogWhenGomaxecsDebugEnvNotEnabled(t *testing.T) {
	tableTest := []struct {
		env string
	}{
		{env: "false"},
		{env: "0"},
		{env: ""},
	}

	for _, tt := range tableTest {
		t.Run("GOMAXECS_DEBUG="+tt.env, func(t *testing.T) {
			t.Setenv("GOMAXECS_DEBUG", tt.env)

			buf := new(bytes.Buffer)
			logger := log.New(buf, "", 0)

			cfg := config.New(config.WithLogger(logger.Printf))

			cfg.DebugLogf("debug log: %s", "stub-message")

			assert.Equal(t, buf.Len(), 0)
		})
	}
}

type mockOption struct {
	isApplied bool
}

func (m *mockOption) Apply(_ *config.Config) {
	m.isApplied = true
}
