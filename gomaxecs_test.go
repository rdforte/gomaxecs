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

// NOTE: This file is intentionally testing a private function.
// This is to ensure that the function runSetMaxProcs is tested
// and remains private to the package.

package gomaxecs //nolint:testpackage // Test private function.

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rdforte/gomaxecs/internal/task/tasktest"
)

func TestGomaxecs_runSetMaxProcs_ECSEnvNotDetected(t *testing.T) {
	t.Parallel()

	curMaxProcs := runtime.GOMAXPROCS(0)

	runSetMaxProcs()

	assert.Equal(t, curMaxProcs, runtime.GOMAXPROCS(0))
}

func TestGomaxecs_runSetMaxProcs_ECSEnvDetected(t *testing.T) {
	curMaxProcs := 1
	runtime.GOMAXPROCS(curMaxProcs)

	wantCPUs := 2
	containerCPU, taskCPU := wantCPUs<<10, wantCPUs

	agent := tasktest.NewECSAgent(t).
		WithContainerMetaEndpoint(containerCPU).
		WithTaskMetaEndpoint(containerCPU, taskCPU).
		Start().
		SetMetaURIEnv()
	defer agent.Close()

	runSetMaxProcs()

	assert.Equal(t, wantCPUs, runtime.GOMAXPROCS(0))
}
