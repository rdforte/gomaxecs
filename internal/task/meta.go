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
	"net/http"
)

// TaskMeta represents the ECS Task Metadata.
type TaskMeta struct {
	Containers []Container `json:"Containers"`
	Limits     Limit       `json:"Limits"` // this is optional in the response
}

// Container represents the ECS Container Metadata.
type Container struct {
	DockerID string `json:"DockerId"` //nolint:tagliatelle // ECS Agent inconsistency. All fields adhere to goPascal but this one.
	Limits   Limit  `json:"Limits"`
}

// Limit contains the CPU limit.
type Limit struct {
	CPU float64 `json:"CPU"`
}

// Grab the container metadata from the ECS Metadata endpoint.
// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v4-examples.html
func (t *Task) getContainerMeta() (Container, error) {
	var container Container

	resp, err := t.client.Get(t.containerMetadataURI)
	if err != nil {
		return container, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return container, newStatusError(resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return container, fmt.Errorf("read failed: %w", err)
	}

	err = json.Unmarshal(data, &container)
	if err != nil {
		return container, fmt.Errorf("unmarshal failed: %w", err)
	}

	return container, nil
}

// Grab the task metadata from the ECS Metadata endpoint + `/task`
// This will also include the container metadata.
// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v4-examples.html#task-metadata-endpoint-v4-example-task-metadata-response
func (t *Task) getTaskMeta() (TaskMeta, error) {
	var task TaskMeta

	resp, err := t.client.Get(t.taskMetadataURI)
	if err != nil {
		return task, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return task, newStatusError(resp.StatusCode)
	}

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

func newStatusError(status int) error {
	return &statusError{status}
}

type statusError struct {
	status int
}

func (e *statusError) Error() string {
	return fmt.Sprintf("request failed, status code: %d", e.status)
}
