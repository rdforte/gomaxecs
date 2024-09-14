package gomaxprocs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getECSTaskMeta(ecsMetaURI string) (Task, error) {
	var task Task

	url := fmt.Sprintf("%s/task", ecsMetaURI)
	resp, err := http.Get(url)
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

type Task struct {
	Containers []Container `json:"Containers"`
	Limits     Limit       `json:"Limits"`
}

type Container struct {
	DockerID string `json:"DockerId"`
	Limits   Limit  `json:"Limits"`
}

type Limit struct {
	CPU float64 `json:"CPU"`
}
