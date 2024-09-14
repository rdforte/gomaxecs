package task_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rdforte/gomax-ecs/internal/task"
)

func TestTask_GetCPU_GetsCPUUsingContainerLimit(t *testing.T) {
	t.Parallel()
	tableTest := []struct {
		name         string
		wantCPU      int
		containerCPU int
		taskCPU      int
		testServer   func(containerCPU, taskCPU int) *httptest.Server
	}{
		{
			name:         "should get cpu of 1 when task CPU limit is 1 and container CPU limit is 1024 vCPU",
			wantCPU:      1,
			containerCPU: 1 << 10,
			taskCPU:      1,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 1 when task CPU limit is 2 and container CPU limit is 1024 vCPU",
			wantCPU:      1,
			containerCPU: 1 << 10,
			taskCPU:      2,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 1 when task CPU limit is 4 and container CPU limit is 1024 vCPU",
			wantCPU:      1,
			containerCPU: 1 << 10,
			taskCPU:      4,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 1 when task CPU limit is 8 and container CPU limit is 1024 vCPU",
			wantCPU:      1,
			containerCPU: 1 << 10,
			taskCPU:      8,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 1 when task CPU limit is 16 and container CPU limit is 1024 vCPU",
			wantCPU:      1,
			containerCPU: 1 << 10,
			taskCPU:      16,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 2 when task CPU limit is 2 and container CPU limit is 2048 vCPU",
			wantCPU:      2,
			containerCPU: 2 << 10,
			taskCPU:      2,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 2 when task CPU limit is 4 and container CPU limit is 2048 vCPU",
			wantCPU:      2,
			containerCPU: 2 << 10,
			taskCPU:      2,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 2 when task CPU limit is 8 and container CPU limit is 2048 vCPU",
			wantCPU:      2,
			containerCPU: 2 << 10,
			taskCPU:      8,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 2 when task CPU limit is 16 and container CPU limit is 2048 vCPU",
			wantCPU:      2,
			containerCPU: 2 << 10,
			taskCPU:      16,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 4 when task CPU limit is 4 and container CPU limit is 4096 vCPU",
			wantCPU:      4,
			containerCPU: 4 << 10,
			taskCPU:      4,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 4 when task CPU limit is 8 and container CPU limit is 4096 vCPU",
			wantCPU:      4,
			containerCPU: 4 << 10,
			taskCPU:      8,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 4 when task CPU limit is 16 and container CPU limit is 4096 vCPU",
			wantCPU:      4,
			containerCPU: 4 << 10,
			taskCPU:      16,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 8 when task CPU limit is 8 and container CPU limit is 8192 vCPU",
			wantCPU:      8,
			containerCPU: 8 << 10,
			taskCPU:      8,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 8 when task CPU limit is 16 and container CPU limit is 8192 vCPU",
			wantCPU:      8,
			containerCPU: 8 << 10,
			taskCPU:      16,
			testServer:   testServerContainerLimit,
		},
		{
			name:         "should get cpu of 16 when task CPU limit is 16 and container CPU limit is 16384 vCPU",
			wantCPU:      16,
			containerCPU: 16 << 10,
			taskCPU:      16,
			testServer:   testServerContainerLimit,
		},
	}

	for _, tt := range tableTest {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ts := tt.testServer(tt.containerCPU, tt.taskCPU)
			defer ts.Close()

			gotCPU, err := task.GetMaxThreads(ts.URL, "container-id")
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCPU, gotCPU)
		})
	}
}

func TestTask_GetCPU_GetsCPUUsingTaskLimit(t *testing.T) {
	t.Parallel()
	tableTest := []struct {
		name       string
		wantCPU    int
		taskCPU    int
		testServer func(taskCPU int) *httptest.Server
	}{
		{
			name:       "should get cpu of 1 when task CPU limit is 0.25 and no container CPU limit set",
			wantCPU:    1,
			taskCPU:    1,
			testServer: testServerTaskLimit,
		},
		{
			name:       "should get cpu of 1 when task CPU limit is 0.5 and no container CPU limit set",
			wantCPU:    1,
			taskCPU:    1,
			testServer: testServerTaskLimit,
		},
		{
			name:       "should get cpu of 1 when task CPU limit is 1 and no container CPU limit set",
			wantCPU:    1,
			taskCPU:    1,
			testServer: testServerTaskLimit,
		},
		{
			name:       "should get cpu of 2 when task CPU limit is 2 and no container CPU limit set",
			wantCPU:    2,
			taskCPU:    2,
			testServer: testServerTaskLimit,
		},
		{
			name:       "should get cpu of 4 when task CPU limit is 4 and no container CPU limit set",
			wantCPU:    4,
			taskCPU:    4,
			testServer: testServerTaskLimit,
		},
		{
			name:       "should get cpu of 8 when task CPU limit is 8 and no container CPU limit set",
			wantCPU:    8,
			taskCPU:    8,
			testServer: testServerTaskLimit,
		},
		{
			name:       "should get cpu of 16 when task CPU limit is 16 and no container CPU limit set",
			wantCPU:    16,
			taskCPU:    16,
			testServer: testServerTaskLimit,
		},
	}

	for _, tt := range tableTest {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ts := tt.testServer(tt.taskCPU)
			defer ts.Close()

			gotCPU, err := task.GetMaxThreads(ts.URL, "container-id")
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCPU, gotCPU)
		})
	}
}

func testServerContainerLimit(containerCPU, taskCPU int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"Containers":[{"DockerId":"container-id","Limits":{"CPU":%d}}],"Limits":{"CPU":%d}}`, containerCPU, taskCPU)))
	}))
}

func testServerTaskLimit(taskCPU int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d}}`, taskCPU)))
	}))
}
