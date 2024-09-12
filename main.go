package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func init() {
	containerURI := os.Getenv("ECS_CONTAINER_METADATA_URI_V4")

	if len(containerURI) == 0 {
		fmt.Println("no container URI")
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/task", containerURI))
	if err != nil {
		fmt.Println("request failed")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read failed")
		log.Fatal(err)
	}

	var task Task
	err = json.Unmarshal(data, &task)
	if err != nil {
		fmt.Println("unmarshal failed")
		log.Fatal(err)
	}

	pathParts := strings.Split(containerURI, "/")
	containerID := pathParts[len(pathParts)-1]

	var containerCPULimit float64
	for _, container := range task.Containers {
		if container.DockerID == containerID {
			fmt.Println("found container")
			containerCPULimit = container.Limits.CPU
		}
	}

	cpu := int(containerCPULimit) >> 10
	proc := runtime.GOMAXPROCS(cpu)
	fmt.Println("GOMAXPROCS", proc)
}

func main() {
	fmt.Println("Hello, World!")
	// fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))
	// fmt.Println("NumCPU", runtime.NumCPU())
	// fmt.Println("Hello, World!")

	// containerURI := os.Getenv("ECS_CONTAINER_METADATA_URI_V4")
	// fmt.Println(containerURI)

	// metaEnabled := os.Getenv("ECS_ENABLE_CONTAINER_METADATA")
	// fmt.Println("enabled", metaEnabled)
	// metaFile := os.Getenv("ECS_CONTAINER_METADATA_FILE")
	// fmt.Println("-->", metaFile)

	// data, err := os.ReadFile(metaFile)
	// if err != nil {
	// fmt.Println("read file failed")
	// log.Fatal(err)
	// }

	// fmt.Println(string(data))

	// if len(containerURI) == 0 {
	// fmt.Println("no container URI")
	// return
	// }

	// resp, err := http.Get(fmt.Sprintf("%s/task", containerURI))
	// if err != nil {
	// fmt.Println("request failed")
	// log.Fatal(err)
	// }
	// defer resp.Body.Close()
	// fmt.Println(resp.Status)
	// data, err := io.ReadAll(resp.Body)
	// if err != nil {
	// fmt.Println("read failed")
	// log.Fatal(err)
	// }

	// var task Task
	// err = json.Unmarshal(data, &task)
	// if err != nil {
	// fmt.Println("unmarshal failed")
	// log.Fatal(err)
	// }

	// fmt.Println(task)

	// pathParts := strings.Split(containerURI, "/")
	// containerID := pathParts[len(pathParts)-1]

	// var containerCPULimit float64
	// for _, container := range task.Containers {
	// if container.DockerID == containerID {
	// fmt.Println("found container")
	// containerCPULimit = container.Limits.CPU
	// }
	// }

	// fmt.Println("containerCPULimit", containerCPULimit)
	// fmt.Println("containerCPULimit", int(containerCPULimit)>>10)

	// entries, err := os.ReadDir("/sys/fs/cgroup/cpu,cpuacct")
	// if err != nil {
	// log.Fatal(err)
	// }

	// for _, entry := range entries {
	// fmt.Println(entry.Name())
	// }
	// fmt.Println("------------------")
	// data, err := os.ReadFile("/sys/fs/cgroup/cpu,cpuacct/cpu.cfs_quota_us")
	// if err != nil {
	// log.Fatal(err)
	// }

	// fmt.Println(string(data))
}

type Task struct {
	Containers []Container `json:"Containers"`
}

type Container struct {
	DockerID string `json:"DockerId"`
	Limits   Limit  `json:"Limits"`
}

type Limit struct {
	CPU float64 `json:"CPU"`
}
