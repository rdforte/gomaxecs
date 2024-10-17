package maxprocs_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rdforte/gomaxecs/maxprocs"
)

const (
	metaURIEnv   = "ECS_CONTAINER_METADATA_URI_V4"
	containerID  = "container-id"
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
	ts := testServerContainerLimit(2<<10, 8)
	defer ts.Close()

	t.Setenv(metaURIEnv, strings.Join([]string{ts.URL, "/", containerID}, ""))

	maxprocs.Set(log.New(io.Discard, "", 0))

	procs := runtime.GOMAXPROCS(0)
	wantProcs := 2
	assert.Equal(t, wantProcs, procs)
}

func TestMaxProcs_Set_LoggerShouldLog(t *testing.T) {
	tableTest := []struct {
		name    string
		wantLog string
		metaURI func() string
	}{
		{
			name:    "should log GOMAXPROCS value when successfully set",
			wantLog: "GOMAXPROCS set to: 2",
			metaURI: func() string {
				ts := testServerContainerLimit(2<<10, 8)
				t.Cleanup(ts.Close)
				return strings.Join([]string{ts.URL, "/", containerID}, "")
			},
		},
		{
			name:    "should log error when task initialisation fails",
			wantLog: "task initialised failed. Unable to set GOMAXPROCS",
			metaURI: func() string {
				return ""
			},
		},
		{
			name:    "should log error when fail to get max procs",
			wantLog: "failed to set GOMAXPROC",
			metaURI: func() string {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
				t.Cleanup(ts.Close)

				return strings.Join([]string{ts.URL, "/", containerID}, "")
			},
		},
	}

	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(metaURIEnv, tt.metaURI())

			buf := new(bytes.Buffer)
			maxprocs.Set(log.New(buf, "", 0))

			assert.Contains(t, buf.String(), tt.wantLog)
		})
	}
}

func testServerContainerLimit(containerCPU, taskCPU int) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/container-id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"DockerId":"container-id"}`, containerCPU)))
	})
	mux.HandleFunc("/container-id/task", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"Containers":[{"DockerId":"container-id","Limits":{"CPU":%d}}],"Limits":{"CPU":%d}}`, containerCPU, taskCPU)))
	})
	return httptest.NewServer(mux)
}
