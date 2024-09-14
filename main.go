package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/rdforte/gomax-ecs/internal/config"
	"github.com/rdforte/gomax-ecs/internal/task"
)

func init() {
	cfg := config.New()
	threads, err := task.GetMaxThreads(cfg.MetadataURI, cfg.ConainerID)
	if err != nil {
		log.Println("failed to set GOMAXPROCS:", err)
		return
	}

	runtime.GOMAXPROCS(threads)
	log.Println("GOMAXPROCS set to:", threads)
}

func main() {
	fmt.Println("main function")
}
