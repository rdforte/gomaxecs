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
