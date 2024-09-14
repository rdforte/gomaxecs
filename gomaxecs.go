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
