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

package task

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/rdforte/gomaxecs/internal/client"
)

// TaskMeta represents the ECS Task Metadata.
type TaskMeta struct {
	Containers []Container `json:"Containers"`
	Limits     Limit       `json:"Limits"`
}

// Container represents the ECS Container Metadata.
type Container struct {
	DockerID string `json:"DockerId"`
	Limits   Limit  `json:"Limits"`
}

// Limit contains the CPU limit.
type Limit struct {
	CPU float64 `json:"CPU"`
}

func (t *Task) getTaskMeta() (TaskMeta, error) {
	var task TaskMeta

	url := fmt.Sprintf("%s/task", t.cfg.MetadataURI)
	client := client.New(t.cfg.Client)
	resp, err := client.Get(url)
	if err != nil {
		return task, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return task, fmt.Errorf("read failed: %w", err)
	}

	err = json.Unmarshal(data, &task)
	if err != nil {
		return task, fmt.Errorf("unmarshal failed: %w", err)
	}

	return task, nil
}
