package config_test

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rdforte/gomaxecs/internal/config"
)

const (
	metaURIEnvKey = "ECS_CONTAINER_METADATA_URI_V4"
	metaURIEnvVal = "mock-ecs-metadata-uri/"
)

func TestConfig_New_LoadConfiguration(t *testing.T) {
	t.Setenv(metaURIEnvKey, metaURIEnvVal)

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
	t.Setenv(metaURIEnvKey, metaURIEnvVal)

	got := config.GetECSMetadataURI()

	want := strings.TrimSuffix(metaURIEnvVal, "/")
	assert.Equal(t, want, got)
}

func TestConfig_GomaxecsDebugEnv_SetsDebugModeInConfig(t *testing.T) {
	t.Setenv(metaURIEnvKey, metaURIEnvVal)

	got := config.GetECSMetadataURI()

	want := strings.TrimSuffix(metaURIEnvVal, "/")
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

			gotLen, wantLen := buf.Len(), 0
			assert.Equal(t, wantLen, gotLen)
		})
	}
}

// write test that checks that default log.Printf is used when debug enabled and no logger set.
func TestConfig_DebugLogf_UsesDefaultLoggerWhenDebugEnabledAndNoLoggerSet(t *testing.T) {
	t.Setenv("GOMAXECS_DEBUG", "true")

	buf := new(bytes.Buffer)

	// replace the original log.Printf with one that writes to our buffer
	originalLogger := log.Default()
	originalOutput := originalLogger.Writer()

	log.SetOutput(buf)
	defer log.SetOutput(originalOutput)

	cfg := config.New()

	// clear buffer so we only have fresh logs
	buf.Reset()

	cfg.DebugLogf("debug log: %s", "stub-message")

	wantLog := "gomaxecs debug: debug log: stub-message\n"
	assert.Contains(t, buf.String(), wantLog)
}

type mockOption struct {
	isApplied bool
}

func (m *mockOption) Apply(_ *config.Config) {
	m.isApplied = true
}
