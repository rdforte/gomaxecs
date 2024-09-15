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

// Package gomaxecs provides a simple way to set GOMAXPROCS based on ECS container
// and task CPU limits.
package gomaxecs

import (
	"log"
	"runtime"

	"github.com/rdforte/gomaxecs/internal/config"
	"github.com/rdforte/gomaxecs/internal/task"
)

func init() {
	cfg := config.New()
	t, err := task.New(cfg)
	if err != nil {
		log.Println("task initialised failed. Unable to set GOMAXPROCS:", err)
	}

	procs, err := t.GetMaxProcs()
	if err != nil {
		log.Println("failed to set GOMAXPROCS:", err)
		return
	}

	runtime.GOMAXPROCS(procs)
	log.Println("GOMAXPROCS set to:", procs)
}
