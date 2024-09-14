package main

import (
	"fmt"
	"log"

	"github.com/rdforte/gomax-ecs/internal/config"
	"github.com/rdforte/gomax-ecs/internal/gomaxprocs"
)

func init() {
	cfg := config.New()
	if cpu, err := gomaxprocs.Set(cfg.MetadataURI, cfg.ConainerID); err != nil {
		log.Println("failed to set GOMAXPROCS:", err)
	} else {
		log.Println("GOMAXPROCS set to:", cpu)
	}
}

func main() {
	fmt.Println("main function")
}
