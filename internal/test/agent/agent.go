package agent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const metaURIEnv = "ECS_CONTAINER_METADATA_URI_V4"

// ECSAgentV4 is a test server that simulates the ECS Agent metadata API.
type ECSAgentV4 struct {
	t            *testing.T
	mux          *http.ServeMux
	server       *httptest.Server
	containerCPU int
}

// NewBuilder builds a new test server that simulates the ECS Agent metadata API.
// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v4.html
func NewV4Builder(t *testing.T) *ECSAgentV4 {
	t.Helper()

	mux := http.NewServeMux()
	return &ECSAgentV4{t, mux, nil, 0}
}

// WithContainerMetaEndpoint sets up the container CPU endpoint on the test server.
func (e *ECSAgentV4) WithContainerMetaEndpoint(containerCPU int) *ECSAgentV4 {
	e.t.Helper()
	e.mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"DockerId":"container-id"}`, containerCPU)))
		assert.NoError(e.t, err)
	})
	return e
}

// WithTaskMetaEndpoint sets up the task CPU endpoint on the test server.
func (e *ECSAgentV4) WithTaskMetaEndpoint(containerCPU, taskCPU int) *ECSAgentV4 {
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

// Start starts the test server.
func (e *ECSAgentV4) Start() *ECSAgentV4 {
	e.t.Helper()
	e.server = httptest.NewServer(e.mux)
	return e
}

// SetMetaURIEnv is a helper function to set the server url for ECS_CONTAINER_METADATA_URI_V4 environment variable.
// This is useful for testing the ECS metadata API.
func (e *ECSAgentV4) SetMetaURIEnv() *ECSAgentV4 {
	e.t.Helper()
	assert.NotNil(e.t, e.server)
	e.t.Setenv(metaURIEnv, e.server.URL)
	return e
}

// Close closes the test server.
func (e *ECSAgentV4) Close() {
	e.server.Close()
}
