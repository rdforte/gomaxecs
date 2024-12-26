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

package gomaxecs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGomaxecs_runSetMaxProcs_NoECSEnvDetected(t *testing.T) {
	curMaxProcs := runtime.GOMAXPROCS(0)
	runSetMaxProcs()
	assert.Equal(t, curMaxProcs, runtime.GOMAXPROCS(0))
}

func TestGomaxecs_runSetMaxProcs_ECSEnvDetected(t *testing.T) {
	curMaxProcs := 1
	runtime.GOMAXPROCS(curMaxProcs)
	// set env variable to simulate ECS environment
	ts := testServerContainerLimit(t, 2<<10, 2)
	t.Setenv("ECS_CONTAINER_METADATA_URI_V4", ts.URL)
	runSetMaxProcs()
	assert.Equal(t, 2, runtime.GOMAXPROCS(0))
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
