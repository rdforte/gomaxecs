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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rdforte/gomaxecs/internal/config"
	"github.com/rdforte/gomaxecs/internal/task"
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
		{
			name:         "should get cpu of 16 when task CPU limit is 0 and container CPU limit is 16384 vCPU",
			wantCPU:      16,
			containerCPU: 16 << 10,
			taskCPU:      0,
			testServer:   testServerContainerLimit,
		},
	}

	for _, tt := range tableTest {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ts := tt.testServer(tt.containerCPU, tt.taskCPU)
			defer ts.Close()

			ecsTask, err := task.New(config.Config{MetadataURI: ts.URL})
			assert.NoError(t, err)

			gotCPU, err := ecsTask.GetMaxProcs()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCPU, gotCPU)
		})
	}
}

func TestTask_GetCPU_GetsCPUUsingContainerLimit_ReturnsError(t *testing.T) {
	t.Parallel()
	tableTest := []struct {
		name         string
		wantError    string
		containerCPU int
		taskCPU      int
		testServer   func(containerCPU, taskCPU int) *httptest.Server
	}{
		{
			name:         "should raise error when task CPU limit is 0 and container CPU limit is 0",
			wantError:    "no CPU limit found for task or container",
			containerCPU: 0,
			taskCPU:      0,
			testServer:   testServerContainerLimit,
		},
	}

	for _, tt := range tableTest {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ts := tt.testServer(tt.containerCPU, tt.taskCPU)
			defer ts.Close()

			ecsTask, err := task.New(config.Config{MetadataURI: ts.URL})
			assert.NoError(t, err)

			_, err = ecsTask.GetMaxProcs()
			assert.ErrorContains(t, err, tt.wantError)
		})
	}
}

func TestTask_GetCPU_GetsCPUUsingContainerLimit_ReturnsErrorWhenContainerEndpointError(t *testing.T) {
	t.Parallel()
	tableTest := []struct {
		name       string
		wantError  string
		testServer func() *httptest.Server
	}{
		{
			name:       "should raise error when ECS container endpoint returns an error",
			wantError:  "failed to get ECS container meta:",
			testServer: testServerContainerEndpointError,
		},
	}

	for _, tt := range tableTest {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ts := tt.testServer()
			defer ts.Close()

			ecsTask, err := task.New(config.Config{MetadataURI: ts.URL})
			assert.NoError(t, err)

			_, err = ecsTask.GetMaxProcs()
			assert.ErrorContains(t, err, tt.wantError)
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

			ecsTask, err := task.New(config.Config{MetadataURI: ts.URL})
			assert.NoError(t, err)

			gotCPU, err := ecsTask.GetMaxProcs()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCPU, gotCPU)
		})
	}
}

func TestTask_GetCPU_GetsCPUUsingTaskLimit_ReturnsError(t *testing.T) {
	t.Parallel()
	tableTest := []struct {
		name       string
		wantError  string
		taskCPU    int
		testServer func(taskCPU int) *httptest.Server
	}{
		{
			name:       "should raise error when task CPU limit is 0 and container CPU limit does not exist",
			wantError:  "no CPU limit found for task or container",
			taskCPU:    0,
			testServer: testServerTaskLimit,
		},
	}

	for _, tt := range tableTest {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ts := tt.testServer(tt.taskCPU)
			defer ts.Close()

			ecsTask, err := task.New(config.Config{MetadataURI: ts.URL})
			assert.NoError(t, err)

			_, err = ecsTask.GetMaxProcs()
			assert.ErrorContains(t, err, tt.wantError)
		})
	}
}

func TestTask_Endpoint_ReturnsError(t *testing.T) {
	t.Parallel()
	tableTest := []struct {
		name       string
		wantError  string
		testServer func() *httptest.Server
	}{
		{
			name:       "should raise error when ECS task endpoint returns an error",
			wantError:  "failed to get ECS task meta:",
			testServer: testServerTaskEndpointError,
		},
	}

	for _, tt := range tableTest {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ts := tt.testServer()
			defer ts.Close()

			ecsTask, err := task.New(config.Config{MetadataURI: ts.URL})
			assert.NoError(t, err)

			_, err = ecsTask.GetMaxProcs()
			assert.ErrorContains(t, err, tt.wantError)
		})
	}
}

func testServerContainerLimit(containerCPU, taskCPU int) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"DockerId":"container-id"}`, containerCPU)))
	})
	mux.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"Containers":[{"DockerId":"container-id","Limits":{"CPU":%d}}]}`, taskCPU, containerCPU)))
	})
	return httptest.NewServer(mux)
}

func testServerTaskLimit(taskCPU int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"DockerId":"container-id"}`, taskCPU)))
	}))
}

func testServerTaskEndpointError() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"DockerId":"container-id"}`, 0)))
	})
	mux.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	return httptest.NewServer(mux)
}

func testServerContainerEndpointError() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	return httptest.NewServer(mux)
}
