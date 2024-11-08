package maxprocs

import (
	"log"
	"runtime"

	"github.com/rdforte/gomaxecs/internal/config"
	"github.com/rdforte/gomaxecs/internal/task"
)

// Set sets GOMAXPROCS based on the CPU limit of the container and the task.
func Set(log *log.Logger) {
	cfg := config.New()
	t := task.New(cfg)

	procs, err := t.GetMaxProcs()
	if err != nil {
		log.Println("failed to set GOMAXPROCS:", err)
		return
	}

	runtime.GOMAXPROCS(procs)
	log.Println("GOMAXPROCS set to:", procs)
}
