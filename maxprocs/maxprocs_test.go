package maxprocs_test

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rdforte/gomaxecs/internal/task/tasktest"
	"github.com/rdforte/gomaxecs/maxprocs"
)

const (
	metaURIEnv   = "ECS_CONTAINER_METADATA_URI_V4"
	taskCPU      = 8
	containerCPU = 2 << 10
)

func TestMaxProcs_Set_SuccessfullySetsGOMAXPROCS(t *testing.T) {
	agent := tasktest.NewECSAgent(t).
		WithContainerMetaEndpoint(containerCPU).
		WithTaskMetaEndpoint(containerCPU, taskCPU).
		Start().
		SetMetaURIEnv()
	defer agent.Close()

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
			wantLog: "maxprocs: Updated GOMAXPROCS=2",
			setup: func(t *testing.T) {
				t.Helper()

				agent := tasktest.NewECSAgent(t).
					WithContainerMetaEndpoint(containerCPU).
					WithTaskMetaEndpoint(containerCPU, taskCPU).
					Start().
					SetMetaURIEnv()

				t.Cleanup(agent.Close)
			},
		},
		{
			name:    "should log GOMAXPROCS value when task cpu limit successfully set",
			wantLog: "maxprocs: Updated GOMAXPROCS=8",
			setup: func(t *testing.T) {
				t.Helper()

				containerCPU := 0
				agent := tasktest.NewECSAgent(t).
					WithContainerMetaEndpoint(containerCPU).
					WithTaskMetaEndpoint(containerCPU, taskCPU).
					Start().
					SetMetaURIEnv()

				t.Cleanup(agent.Close)
			},
		},
		{
			name:    "should log error when fail to get max procs",
			wantLog: "maxprocs: Failed to set GOMAXPROCS",
			setup: func(t *testing.T) {
				t.Helper()

				agent := tasktest.NewECSAgent(t).
					WithContainerMetaEndpointInternalServerError().
					WithTaskMetaEndpointInternalServerError().
					Start().
					SetMetaURIEnv()

				t.Cleanup(agent.Close)
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

func TestMaxProcs_Set_UndoLogsNoChangesWhenHonorsGOMAXPROCSEnv(t *testing.T) {
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

	agent := tasktest.NewECSAgent(t).
		WithContainerMetaEndpoint(containerCPU).
		WithTaskMetaEndpoint(containerCPU, taskCPU).
		Start().
		SetMetaURIEnv()
	defer agent.Close()

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
