package maxprocs

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/rdforte/gomaxecs/internal/config"
	ecstask "github.com/rdforte/gomaxecs/internal/task"
)

const maxProcsKey = "GOMAXPROCS"

// Set sets GOMAXPROCS based on the CPU limit of the container and the task.
// returns a function to reset GOMAXPROCS to its previous value and an error if one occurred.
func Set(opts ...config.Option) (func(), error) {
	cfg := config.New(opts...)
	task := ecstask.New(cfg)

	undoNoop := func() {
		cfg.Log("maxprocs: No GOMAXPROCS change to reset")
	}

	if curMaxProcs, exists := honorCurrentMaxProcs(); exists {
		cfg.Log("maxprocs: Honoring GOMAXPROCS=%q as set in environment", curMaxProcs)
		return undoNoop, nil
	}

	prevProcs := currentMaxProcs()
	undo := func() {
		cfg.Log("maxprocs: Resetting GOMAXPROCS to %v", prevProcs)
		runtime.GOMAXPROCS(prevProcs)
	}

	procs, err := task.GetMaxProcs(context.Background())
	if err != nil {
		cfg.Log("maxprocs: Failed to set GOMAXPROCS:", err)
		return undo, fmt.Errorf("failed to set GOMAXPROCS: %w", err)
	}

	runtime.GOMAXPROCS(procs)
	cfg.Log("maxprocs: Updating GOMAXPROCS=%v", procs)

	return undo, nil
}

func honorCurrentMaxProcs() (string, bool) {
	return os.LookupEnv(maxProcsKey)
}

func currentMaxProcs() int {
	return runtime.GOMAXPROCS(0)
}

// WithLogger sets the logger. By default, no logger is set.
func WithLogger(printf func(format string, args ...any)) config.Option {
	return config.WithLogger(printf)
}

// IsECS returns true if detected ECS environment.
func IsECS() bool {
	return len(config.GetECSMetadataURI()) > 0
}
