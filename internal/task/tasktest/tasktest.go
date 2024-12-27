package tasktest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	metaURIEnv   = "ECS_CONTAINER_METADATA_URI_V4"
	taskMetaPath = "/task"
)

// ECSAgent is a test server that simulates the ECS Agent metadata API.
type ECSAgent struct {
	t            *testing.T
	mux          *http.ServeMux
	server       *httptest.Server
	containerCPU int
}

// NewECSAgent builds a new test server that simulates the ECS Agent metadata API.
// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v4.html
func NewECSAgent(t *testing.T) *ECSAgent {
	t.Helper()
	mux := http.NewServeMux()
	return &ECSAgent{t, mux, nil, 0}
}

// WithContainerMetaEndpoint sets up the container CPU endpoint on the test server.
func (e *ECSAgent) WithContainerMetaEndpoint(containerCPU int) *ECSAgent {
	e.t.Helper()
	e.mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"DockerId":"container-id"}`, containerCPU)))
		assert.NoError(e.t, err)
	})
	return e
}

// WithTaskMetaEndpoint sets up the task CPU endpoint on the test server.
func (e *ECSAgent) WithTaskMetaEndpoint(containerCPU, taskCPU int) *ECSAgent {
	e.t.Helper()
	e.mux.HandleFunc("/task", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(fmt.Sprintf(
			`{"Containers":[{"DockerId":"container-id","Limits":{"CPU":%d}}],"Limits":{"CPU":%d}}`,
			containerCPU,
			taskCPU,
		)))
		assert.NoError(e.t, err)
	})
	return e
}

// WithContainerMetaEndpointInternalServerError sets up the container meta endpoint to return an internal server error.
func (e *ECSAgent) WithContainerMetaEndpointInternalServerError() *ECSAgent {
	e.t.Helper()
	e.mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	return e
}

// WithTaskMetaEndpointInternalServerError sets up the task meta endpoint to return an internal server error.
func (e *ECSAgent) WithTaskMetaEndpointInternalServerError() *ECSAgent {
	e.t.Helper()
	e.mux.HandleFunc("/task", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	return e
}

// Start starts the test server.
func (e *ECSAgent) Start() *ECSAgent {
	e.t.Helper()
	e.server = httptest.NewServer(e.mux)
	return e
}

// SetMetaURIEnv is a helper function to set the server url for ECS_CONTAINER_METADATA_URI_V4 environment variable.
// This is useful for testing the ECS metadata API.
func (e *ECSAgent) SetMetaURIEnv() *ECSAgent {
	e.t.Helper()
	assert.NotNil(e.t, e.server)
	e.t.Setenv(metaURIEnv, e.server.URL)
	return e
}

// Close closes the test server.
func (e *ECSAgent) Close() {
	e.server.Close()
}

// GetContainerMetaEndpoint returns the container metadata endpoint.
func (e *ECSAgent) GetContainerMetaEndpoint() string {
	return e.server.URL
}

// GetTaskMetaEndpoint returns the task metadata endpoint.
func (e *ECSAgent) GetTaskMetaEndpoint() string {
	return e.server.URL + taskMetaPath
}
