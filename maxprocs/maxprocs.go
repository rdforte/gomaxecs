package maxprocs

import (
	"context"
	"os"
	"runtime"

	"github.com/rdforte/gomaxecs/internal/config"
	"github.com/rdforte/gomaxecs/internal/task"
)

const maxProcsKey = "GOMAXPROCS"

// Set sets GOMAXPROCS based on the CPU limit of the container and the task.
// returns a function to reset GOMAXPROCS to its previous value and an error if any.
func Set(opts ...config.Option) (func(), error) {
	cfg := config.New(opts...)
	t := task.New(cfg)

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

	procs, err := t.GetMaxProcs(context.Background())
	if err != nil {
		cfg.Log("maxprocs: Failed to set GOMAXPROCS:", err)
		return undo, err
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
