package maxprocs_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rdforte/gomaxecs/maxprocs"
)

const (
	metaURIEnv   = "ECS_CONTAINER_METADATA_URI_V4"
	taskCPU      = 8
	containerCPU = 2 << 10
)

func TestMain(m *testing.M) {
	if err := os.Unsetenv(metaURIEnv); err != nil {
		log.Fatalf("failed to unset %s: %v", metaURIEnv, err)
	}

	os.Exit(m.Run())
}

func TestMaxProcs_Set_SuccessfullySetsGOMAXPROCS(t *testing.T) {
	ts := testServerContainerLimit(t, containerCPU, taskCPU)
	defer ts.Close()

	t.Setenv(metaURIEnv, ts.URL)

	_, err := maxprocs.Set()
	require.NoError(t, err)

	procs := runtime.GOMAXPROCS(0)
	wantProcs := 2
	assert.Equal(t, wantProcs, procs)
}

func TestMaxProcs_Set_LoggerShouldLog(t *testing.T) {
	tableTest := []struct {
		name    string
		wantLog string
		setup   func(t *testing.T)
	}{
		{
			name:    "should log when honors current max procs",
			wantLog: "maxprocs: Honoring GOMAXPROCS=\"4\" as set in environment",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv("GOMAXPROCS", "4")
			},
		},
		{
			name:    "should log GOMAXPROCS value when container cpu limit successfully set",
			wantLog: "maxprocs: Updating GOMAXPROCS=2",
			setup: func(t *testing.T) {
				t.Helper()

				ts := testServerContainerLimit(t, containerCPU, taskCPU)
				t.Cleanup(ts.Close)

				t.Setenv(metaURIEnv, ts.URL)
			},
		},
		{
			name:    "should log GOMAXPROCS value when task cpu limit successfully set",
			wantLog: "maxprocs: Updating GOMAXPROCS=8",
			setup: func(t *testing.T) {
				t.Helper()

				ts := testServerContainerLimit(t, 0, taskCPU)
				t.Cleanup(ts.Close)

				t.Setenv(metaURIEnv, ts.URL)
			},
		},
		{
			name:    "should log error when fail to get max procs",
			wantLog: "maxprocs: Failed to set GOMAXPROCS",
			setup: func(t *testing.T) {
				t.Helper()

				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
				t.Cleanup(ts.Close)

				t.Setenv(metaURIEnv, ts.URL)
			},
		},
	}

	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			buf := new(bytes.Buffer)
			logger := log.New(buf, "", 0)

			_, _ = maxprocs.Set(maxprocs.WithLogger(logger.Printf))

			assert.Contains(t, buf.String(), tt.wantLog)
		})
	}
}

func TestMaxProcs_Set_UndoLogsNoChangesWhenHonorsCurrentMaxProcs(t *testing.T) {
	t.Setenv("GOMAXPROCS", "4")

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	undo, _ := maxprocs.Set(maxprocs.WithLogger(logger.Printf))

	undo()

	assert.Contains(t, buf.String(), "maxprocs: No GOMAXPROCS change to reset")
}

func TestMaxProcs_Set_UndoResetsGOMAXPROCS(t *testing.T) {
	initialProcs := 5
	runtime.GOMAXPROCS(initialProcs)

	taskCPU := 10
	containerCPU := 0

	ts := testServerContainerLimit(t, containerCPU, taskCPU)
	defer ts.Close()

	t.Setenv(metaURIEnv, ts.URL)

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	undo, _ := maxprocs.Set(maxprocs.WithLogger(logger.Printf))

	assert.Equal(t, taskCPU, runtime.GOMAXPROCS(0)) // GOMAXPROCS should be set to taskCPU

	undo() // reset GOMAXPROCS

	assert.Equal(t, initialProcs, runtime.GOMAXPROCS(0)) // GOMAXPROCS should be reset to initialProcs

	assert.Contains(t, buf.String(), fmt.Sprintf("maxprocs: Resetting GOMAXPROCS to %v", initialProcs))
}

func TestMaxProcs_IsECS_ReturnsTrueIfDetectedECSEnvironment(t *testing.T) {
	t.Setenv(metaURIEnv, "mock-ecs-metadata-uri")
	assert.True(t, maxprocs.IsECS())
}

func TestMaxProcs_IsECS_ReturnsFalseIfNotDetectedECSEnvironment(t *testing.T) {
	assert.False(t, maxprocs.IsECS())
}

func testServerContainerLimit(t *testing.T, containerCPU, taskCPU int) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"DockerId":"container-id"}`, containerCPU)))
		assert.NoError(t, err)
	})
	mux.HandleFunc("/task", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(fmt.Sprintf(
			`{"Containers":[{"DockerId":"container-id","Limits":{"CPU":%d}}],"Limits":{"CPU":%d}}`,
			containerCPU,
			taskCPU,
		)))
		assert.NoError(t, err)
	})

	return httptest.NewServer(mux)
}
