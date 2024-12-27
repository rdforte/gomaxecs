package agent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ECSAgent struct {
	server *httptest.Server
	t      *testing.T
}

// New is a helper function to create a new test server that simulates
// the ECS Agent metadata API.
// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v4.html
func New(t *testing.T, containerCPU, taskCPU int) *ECSAgent {
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

	return &ECSAgent{httptest.NewServer(mux), t}
}

// SetServerURL is a helper function to set the server url for ECS_CONTAINER_METADATA_URI_V4 environment variable.
// This is useful for testing the ECS metadata API.
func (e *ECSAgent) SetServerURL() {
	e.t.Helper()
	e.t.Setenv("ECS_CONTAINER_METADATA_URI_V4", e.server.URL)
}
