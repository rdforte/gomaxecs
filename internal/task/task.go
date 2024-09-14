package task

import (
	"fmt"

	"github.com/rdforte/gomax-ecs/internal/config"
)

// Task represents a task.
type Task struct {
	cfg config.Config
}

// New returns a new Task.
func New(cfg config.Config) (*Task, error) {
	if len(cfg.MetadataURI) == 0 {
		return nil, fmt.Errorf("no container URI provided")
	}

	if len(cfg.ConainerID) == 0 {
		return nil, fmt.Errorf("no container ID provided")
	}
	return &Task{cfg}, nil
}

// GetMaxProcs is responsible for getting the max number of processors, or
// /sched/gomaxprocs:threads based on the CPU limit of the container and the task.
// The container vCPU can not be greater than Task CPU limit, therefore if
// Task CPU limit is less than 1, the max threads returned is 1.
// If no CPU limit is found for the container, then the max number of threads
// returned is the number of CPU's for the ECS Task. If no CPU limit is found for the
// Task, then 0 is returned with an error
func (t *Task) GetMaxProcs() (int, error) {
	task, err := t.getTaskMeta()
	if err != nil {
		return 0, fmt.Errorf("failed to get ECS task: %w", err)
	}

	if task.Limits.CPU == 0 {
		return 0, fmt.Errorf("no CPU limit found for task")
	}

	var cpuLimit float64
	for _, container := range task.Containers {
		if container.DockerID == t.cfg.ConainerID {
			cpuLimit = container.Limits.CPU
		}
	}

	if cpuLimit == 0 {
		minThreads := 1
		return max(int(task.Limits.CPU), minThreads), nil
	}

	cpu := int(cpuLimit) >> 10

	return cpu, nil
}
