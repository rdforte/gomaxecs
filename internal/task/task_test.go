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

const taskMetaPath = "/task"

func TestTask_GetMaxProcs_GetsCPUUsingContainerLimit(t *testing.T) {
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

			containerURI, taskURI := buildMetaEndpoints(ts)
			ecsTask := task.New(config.Config{ContainerMetadataURI: containerURI, TaskMetadataURI: taskURI})

			gotCPU, err := ecsTask.GetMaxProcs()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCPU, gotCPU)
		})
	}
}

func TestTask_GetMaxProcs_GetsCPUUsingTaskLimit(t *testing.T) {
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

			containerURI, taskURI := buildMetaEndpoints(ts)
			ecsTask := task.New(config.Config{ContainerMetadataURI: containerURI, TaskMetadataURI: taskURI})

			gotCPU, err := ecsTask.GetMaxProcs()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCPU, gotCPU)
		})
	}
}

func TestTask_GetMaxProcs_ReturnsErrorWhenFailToGetNumCPU(t *testing.T) {
	t.Parallel()

	tableTest := []struct {
		name         string
		wantError    string
		containerCPU int
		taskCPU      int
		testServer   func(t *testing.T, containerCPU, taskCPU int) (containerMetaURI, taskMetaURI string)
	}{
		{
			name:         "should raise error when task CPU limit is 0 and container CPU limit is 0",
			wantError:    "no CPU limit found for task or container",
			containerCPU: 0,
			taskCPU:      0,
			testServer: func(t *testing.T, containerCPU, taskCPU int) (string, string) {
				t.Helper()

				ts := testServerContainerLimit(containerCPU, taskCPU)

				t.Cleanup(func() {
					ts.Close()
				})

				return buildMetaEndpoints(ts)
			},
		},
		{
			name:      "should raise error when task CPU limit is 0 and container CPU limit does not exist",
			wantError: "no CPU limit found for task or container",
			taskCPU:   0,
			testServer: func(t *testing.T, _, taskCPU int) (string, string) {
				t.Helper()

				ts := testServerTaskLimit(taskCPU)

				t.Cleanup(func() {
					ts.Close()
				})

				return buildMetaEndpoints(ts)
			},
		},
		{
			name:      "should raise error when ECS container endpoint is not 200 OK",
			wantError: "failed to get ECS container meta: request failed, status code: 500",
			testServer: func(t *testing.T, _, _ int) (string, string) {
				t.Helper()

				mux := http.NewServeMux()
				mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})
				ts := httptest.NewServer(mux)

				t.Cleanup(func() {
					ts.Close()
				})

				return buildMetaEndpoints(ts)
			},
		},
		{
			name:      "should raise error when ECS task endpoint is not 200 OK",
			wantError: "failed to get ECS task meta: request failed, status code: 500",
			testServer: func(t *testing.T, _, _ int) (string, string) {
				t.Helper()

				mux := http.NewServeMux()
				mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					_, err := w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"DockerId":"container-id"}`, 0)))
					assert.NoError(t, err)
				})
				mux.HandleFunc(taskMetaPath, func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})
				ts := httptest.NewServer(mux)

				t.Cleanup(func() {
					ts.Close()
				})

				return buildMetaEndpoints(ts)
			},
		},
		{
			name:      "should raise error when fail to read ECS container meta",
			wantError: "failed to get ECS container meta: read failed",
			testServer: func(t *testing.T, _, _ int) (string, string) {
				t.Helper()

				mux := http.NewServeMux()
				mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					_, err := w.Write([]byte("partial-data"))
					assert.NoError(t, err)
					conn, _, _ := w.(http.Hijacker).Hijack()
					conn.Close()
				})
				mux.HandleFunc(taskMetaPath, func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"Containers":[{"DockerId":"container-id","Limits":{"CPU":%d}}]}`, 1, 1024)))
				})
				ts := httptest.NewServer(mux)

				t.Cleanup(func() {
					ts.Close()
				})

				return buildMetaEndpoints(ts)
			},
		},
		{
			name:      "should raise error when fail to read ECS task meta",
			wantError: "failed to get ECS task meta: read failed",
			testServer: func(t *testing.T, _, _ int) (string, string) {
				t.Helper()

				mux := http.NewServeMux()
				mux.HandleFunc(taskMetaPath, func(w http.ResponseWriter, r *http.Request) {
					_, err := w.Write([]byte("partial-data"))
					assert.NoError(t, err)
					conn, _, _ := w.(http.Hijacker).Hijack()
					conn.Close()
				})
				mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					_, err := w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"DockerId":"container-id"}`, 0)))
					assert.NoError(t, err)
				})
				ts := httptest.NewServer(mux)

				t.Cleanup(func() {
					ts.Close()
				})

				return buildMetaEndpoints(ts)
			},
		},
		{
			name:      "should raise error when fail to unmarshal ECS container meta",
			wantError: "failed to get ECS container meta: unmarshal failed",
			testServer: func(t *testing.T, _, _ int) (string, string) {
				t.Helper()

				mux := http.NewServeMux()
				mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					_, err := w.Write([]byte("invlaid-json"))
					assert.NoError(t, err)
				})
				mux.HandleFunc(taskMetaPath, func(w http.ResponseWriter, r *http.Request) {
					_, err := w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"Containers":[{"DockerId":"container-id","Limits":{"CPU":%d}}]}`, 1, 1024)))
					assert.NoError(t, err)
				})
				ts := httptest.NewServer(mux)

				t.Cleanup(func() {
					ts.Close()
				})

				return buildMetaEndpoints(ts)
			},
		},
		{
			name:      "should raise error when fail to unmarshal ECS task meta",
			wantError: "failed to get ECS task meta: unmarshal failed",
			testServer: func(t *testing.T, _, _ int) (string, string) {
				t.Helper()

				mux := http.NewServeMux()
				mux.HandleFunc(taskMetaPath, func(w http.ResponseWriter, r *http.Request) {
					_, err := w.Write([]byte("invlaid-json"))
					assert.NoError(t, err)
				})
				mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"DockerId":"container-id"}`, 0)))
				})
				ts := httptest.NewServer(mux)

				t.Cleanup(func() {
					ts.Close()
				})

				return buildMetaEndpoints(ts)
			},
		},
		{
			name:      "should raise error when fail to get ECS container meta",
			wantError: "failed to get ECS container meta: request failed",
			testServer: func(t *testing.T, _, _ int) (string, string) {
				t.Helper()

				cpu := 1
				ts := testServerTaskLimit(cpu)

				t.Cleanup(func() {
					ts.Close()
				})

				_, taskURI := buildMetaEndpoints(ts)
				containerURI := "invalid-uri"

				return containerURI, taskURI
			},
		},
		{
			name:      "should raise error when fail to get ECS task meta",
			wantError: "failed to get ECS task meta: request failed",
			testServer: func(t *testing.T, _, _ int) (string, string) {
				t.Helper()

				cpu := 1
				ts := testServerTaskLimit(cpu)

				t.Cleanup(func() {
					ts.Close()
				})

				containerURI, _ := buildMetaEndpoints(ts)
				taskURI := "invalid-uri"

				return containerURI, taskURI
			},
		},
	}

	for _, tt := range tableTest {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			containerMetaURI, taskMetaURI := tt.testServer(t, tt.containerCPU, tt.taskCPU)

			ecsTask := task.New(config.Config{ContainerMetadataURI: containerMetaURI, TaskMetadataURI: taskMetaURI})

			_, err := ecsTask.GetMaxProcs()
			assert.ErrorContains(t, err, tt.wantError)
		})
	}
}

func testServerContainerLimit(containerCPU, taskCPU int) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"DockerId":"container-id"}`, containerCPU)))
	})
	mux.HandleFunc(taskMetaPath, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"Containers":[{"DockerId":"container-id","Limits":{"CPU":%d}}]}`, taskCPU, containerCPU)))
	})
	return httptest.NewServer(mux)
}

func testServerTaskLimit(taskCPU int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"Limits":{"CPU":%d},"DockerId":"container-id"}`, taskCPU)))
	}))
}

func buildMetaEndpoints(ts *httptest.Server) (containerMetaURI, taskMetaURI string) {
	return ts.URL, ts.URL + taskMetaPath
}
