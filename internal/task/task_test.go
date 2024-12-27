// Copyright 2004 Ryan Forte
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rdforte/gomaxecs/internal/config"
	"github.com/rdforte/gomaxecs/internal/task"
	"github.com/rdforte/gomaxecs/internal/task/tasktest"
)

func TestTask_GetMaxProcs_GetsCPUUsingContainerLimit(t *testing.T) {
	t.Parallel()

	tableTest := []struct {
		name         string
		wantCPU      int
		containerCPU int
		taskCPU      int
	}{
		{
			name:         "should get cpu of 1 when task CPU limit is 1 and container CPU limit is 512 vCPU",
			wantCPU:      1,
			containerCPU: 1 << 9,
			taskCPU:      1,
		},
		{
			name:         "should get cpu of 1 when task CPU limit is 1 and container CPU limit is 1024 vCPU",
			wantCPU:      1,
			containerCPU: 1 << 10,
			taskCPU:      1,
		},
		{
			name:         "should get cpu of 1 when task CPU limit is 2 and container CPU limit is 1024 vCPU",
			wantCPU:      1,
			containerCPU: 1 << 10,
			taskCPU:      2,
		},
		{
			name:         "should get cpu of 1 when task CPU limit is 4 and container CPU limit is 1024 vCPU",
			wantCPU:      1,
			containerCPU: 1 << 10,
			taskCPU:      4,
		},
		{
			name:         "should get cpu of 1 when task CPU limit is 8 and container CPU limit is 1024 vCPU",
			wantCPU:      1,
			containerCPU: 1 << 10,
			taskCPU:      8,
		},
		{
			name:         "should get cpu of 1 when task CPU limit is 16 and container CPU limit is 1024 vCPU",
			wantCPU:      1,
			containerCPU: 1 << 10,
			taskCPU:      16,
		},
		{
			name:         "should get cpu of 2 when task CPU limit is 2 and container CPU limit is 2048 vCPU",
			wantCPU:      2,
			containerCPU: 2 << 10,
			taskCPU:      2,
		},
		{
			name:         "should get cpu of 2 when task CPU limit is 4 and container CPU limit is 2048 vCPU",
			wantCPU:      2,
			containerCPU: 2 << 10,
			taskCPU:      2,
		},
		{
			name:         "should get cpu of 2 when task CPU limit is 8 and container CPU limit is 2048 vCPU",
			wantCPU:      2,
			containerCPU: 2 << 10,
			taskCPU:      8,
		},
		{
			name:         "should get cpu of 2 when task CPU limit is 16 and container CPU limit is 2048 vCPU",
			wantCPU:      2,
			containerCPU: 2 << 10,
			taskCPU:      16,
		},
		{
			name:         "should get cpu of 4 when task CPU limit is 4 and container CPU limit is 4096 vCPU",
			wantCPU:      4,
			containerCPU: 4 << 10,
			taskCPU:      4,
		},
		{
			name:         "should get cpu of 4 when task CPU limit is 8 and container CPU limit is 4096 vCPU",
			wantCPU:      4,
			containerCPU: 4 << 10,
			taskCPU:      8,
		},
		{
			name:         "should get cpu of 4 when task CPU limit is 16 and container CPU limit is 4096 vCPU",
			wantCPU:      4,
			containerCPU: 4 << 10,
			taskCPU:      16,
		},
		{
			name:         "should get cpu of 8 when task CPU limit is 8 and container CPU limit is 8192 vCPU",
			wantCPU:      8,
			containerCPU: 8 << 10,
			taskCPU:      8,
		},
		{
			name:         "should get cpu of 8 when task CPU limit is 16 and container CPU limit is 8192 vCPU",
			wantCPU:      8,
			containerCPU: 8 << 10,
			taskCPU:      16,
		},
		{
			name:         "should get cpu of 16 when task CPU limit is 16 and container CPU limit is 16384 vCPU",
			wantCPU:      16,
			containerCPU: 16 << 10,
			taskCPU:      16,
		},
		// For tasks that are hosted on Amazon EC2 instances, the CPU limit is optional.
		// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html#task_size
		{
			name:         "should get cpu of 16 when task CPU limit is 0 and container CPU limit is 16384 vCPU",
			wantCPU:      16,
			containerCPU: 16 << 10,
			taskCPU:      0,
		},
	}

	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			agent := tasktest.NewECSAgent(t).
				WithContainerMetaEndpoint(tt.containerCPU).
				WithTaskMetaEndpoint(tt.containerCPU, tt.taskCPU).
				Start()
			defer agent.Close()

			ecsTask := task.New(config.Config{
				ContainerMetadataURI: agent.GetContainerMetaEndpoint(),
				TaskMetadataURI:      agent.GetTaskMetaEndpoint(),
			})

			gotCPU, err := ecsTask.GetMaxProcs(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tt.wantCPU, gotCPU)
		})
	}
}

func TestTask_GetMaxProcs_GetsCPUUsingTaskLimit(t *testing.T) {
	t.Parallel()

	tableTest := []struct {
		name    string
		wantCPU int
		taskCPU int
	}{
		{
			name:    "should get cpu of 1 when task CPU limit is 0.25 and no container CPU limit set",
			wantCPU: 1,
			taskCPU: 1,
		},
		{
			name:    "should get cpu of 1 when task CPU limit is 0.5 and no container CPU limit set",
			wantCPU: 1,
			taskCPU: 1,
		},
		{
			name:    "should get cpu of 1 when task CPU limit is 1 and no container CPU limit set",
			wantCPU: 1,
			taskCPU: 1,
		},
		{
			name:    "should get cpu of 2 when task CPU limit is 2 and no container CPU limit set",
			wantCPU: 2,
			taskCPU: 2,
		},
		{
			name:    "should get cpu of 4 when task CPU limit is 4 and no container CPU limit set",
			wantCPU: 4,
			taskCPU: 4,
		},
		{
			name:    "should get cpu of 8 when task CPU limit is 8 and no container CPU limit set",
			wantCPU: 8,
			taskCPU: 8,
		},
		{
			name:    "should get cpu of 16 when task CPU limit is 16 and no container CPU limit set",
			wantCPU: 16,
			taskCPU: 16,
		},
	}

	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			containerCPU := 0

			agent := tasktest.NewECSAgent(t).
				WithContainerMetaEndpoint(containerCPU).
				WithTaskMetaEndpoint(containerCPU, tt.taskCPU).
				Start()
			defer agent.Close()

			ecsTask := task.New(config.Config{
				ContainerMetadataURI: agent.GetContainerMetaEndpoint(),
				TaskMetadataURI:      agent.GetTaskMetaEndpoint(),
			})

			gotCPU, err := ecsTask.GetMaxProcs(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tt.wantCPU, gotCPU)
		})
	}
}

func TestTask_GetMaxProcs_ReturnsErrorWhenFailToGetNumCPU(t *testing.T) {
	t.Parallel()

	tableTest := []struct {
		name       string
		wantError  string
		testServer func(t *testing.T) (containerMetaURI, taskMetaURI string)
	}{
		{
			name:      "should raise error when task CPU limit is 0 and container CPU limit is 0",
			wantError: "no CPU limit found for task or container",
			testServer: func(t *testing.T) (string, string) {
				t.Helper()

				containerCPU, taskCPU := 0, 0

				agent := tasktest.NewECSAgent(t).
					WithContainerMetaEndpoint(containerCPU).
					WithTaskMetaEndpoint(containerCPU, taskCPU).
					Start()

				t.Cleanup(agent.Close)

				return agent.GetContainerMetaEndpoint(), agent.GetTaskMetaEndpoint()
			},
		},
		{
			name:      "should raise error when ECS container endpoint is not 200 OK",
			wantError: "failed to get ECS container meta: request failed, status code: 500",
			testServer: func(t *testing.T) (string, string) {
				t.Helper()

				containerCPU, taskCPU := 0, 1
				agent := tasktest.NewECSAgent(t).
					WithTaskMetaEndpoint(containerCPU, taskCPU).
					WithContainerMetaEndpointInternalServerError().
					Start()

				t.Cleanup(agent.Close)

				return agent.GetContainerMetaEndpoint(), agent.GetTaskMetaEndpoint()
			},
		},
		{
			name:      "should raise error when ECS task endpoint is not 200 OK",
			wantError: "failed to get ECS task meta: request failed, status code: 500",
			testServer: func(t *testing.T) (string, string) {
				t.Helper()

				containerCPU := 1
				agent := tasktest.NewECSAgent(t).
					WithContainerMetaEndpoint(containerCPU).
					WithTaskMetaEndpointInternalServerError().
					Start()

				t.Cleanup(agent.Close)

				return agent.GetContainerMetaEndpoint(), agent.GetTaskMetaEndpoint()
			},
		},
		{
			name:      "should raise error when fail to unmarshal ECS container meta",
			wantError: "failed to get ECS container meta: unmarshal failed",
			testServer: func(t *testing.T) (string, string) {
				t.Helper()

				containerCPU, taskCPU := 1<<10, 1
				agent := tasktest.NewECSAgent(t).
					WithTaskMetaEndpoint(containerCPU, taskCPU).
					WithContainerMetaEndpointInvalidJSON().
					Start()

				t.Cleanup(agent.Close)

				return agent.GetContainerMetaEndpoint(), agent.GetTaskMetaEndpoint()
			},
		},
		{
			name:      "should raise error when fail to unmarshal ECS task meta",
			wantError: "failed to get ECS task meta: unmarshal failed",
			testServer: func(t *testing.T) (string, string) {
				t.Helper()

				containerCPU := 1 << 10
				agent := tasktest.NewECSAgent(t).
					WithContainerMetaEndpoint(containerCPU).
					WithTaskMetaEndpointInvalidJSON().
					Start()

				t.Cleanup(agent.Close)

				return agent.GetContainerMetaEndpoint(), agent.GetTaskMetaEndpoint()
			},
		},
		{
			name:      "should raise error when fail to get ECS container meta",
			wantError: "failed to get ECS container meta: request failed",
			testServer: func(t *testing.T) (string, string) {
				t.Helper()

				containerCPU, taskCPU := 1<<10, 1
				agent := tasktest.NewECSAgent(t).
					WithContainerMetaEndpoint(containerCPU).
					WithTaskMetaEndpoint(containerCPU, taskCPU).
					Start()

				t.Cleanup(agent.Close)

				containerURI := "invalid-uri"
				taskURI := agent.GetTaskMetaEndpoint()

				return containerURI, taskURI
			},
		},
		{
			name:      "should raise error when fail to get ECS task meta",
			wantError: "failed to get ECS task meta: request failed",
			testServer: func(t *testing.T) (string, string) {
				t.Helper()

				containerCPU, taskCPU := 1<<10, 1
				agent := tasktest.NewECSAgent(t).
					WithContainerMetaEndpoint(containerCPU).
					WithTaskMetaEndpoint(containerCPU, taskCPU).
					Start()

				t.Cleanup(agent.Close)

				containerURI := agent.GetContainerMetaEndpoint()
				taskURI := "invalid-uri"

				return containerURI, taskURI
			},
		},
	}

	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			containerMetaURI, taskMetaURI := tt.testServer(t)
			ecsTask := task.New(config.Config{ContainerMetadataURI: containerMetaURI, TaskMetadataURI: taskMetaURI})

			_, err := ecsTask.GetMaxProcs(context.Background())
			assert.ErrorContains(t, err, tt.wantError)
		})
	}
}
