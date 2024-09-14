package gomaxprocs

import (
	"fmt"
	"runtime"
)

// Set sets the GOMAXPROCS value based on the CPU limit of the container.
// The container vCPU can not be greater than Task CPU limit, therefore if
// Task CPU limit is less than 1, GOMAXPROCS is set to 1.
// If no CPU limit is found for the container, then GOMAXPROCS is not set
// therefore GOMAXPROCS reverts back to runtime.NumCPU() which is the number
// of CPUs set on the ECS Task.
func Set(ecsMetaURI, containerID string) (int, error) {
	defaultCPU := runtime.NumCPU()
	if len(ecsMetaURI) == 0 {
		return defaultCPU, fmt.Errorf("no container URI provided")
	}

	if len(containerID) == 0 {
		return defaultCPU, fmt.Errorf("no container ID provided")
	}

	task, err := getECSTaskMeta(ecsMetaURI)
	if err != nil {
		return defaultCPU, fmt.Errorf("failed to get ECS task: %w", err)
	}

	if task.Limits.CPU < 1 {
		runtime.GOMAXPROCS(1)
		return 1, nil
	}

	var cpuLimit float64
	for _, container := range task.Containers {
		if container.DockerID == containerID {
			cpuLimit = container.Limits.CPU
		}
	}

	if cpuLimit == 0 {
		return defaultCPU, nil
	}

	cpu := int(cpuLimit) >> 10
	runtime.GOMAXPROCS(cpu)

	return cpu, nil
}
