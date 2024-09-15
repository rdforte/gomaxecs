# gomaxecs [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![License: MIT][mit-img]][mit]

Package for auto setting GOMAXPROCS based on ECS task and container CPU limits.

Due to Go not being CFS aware  https://github.com/golang/go/issues/33803 and because [uber automaxprocs](https://github.com/uber-go/automaxprocs) is unable to set GOMAXPROCS for ECS https://github.com/uber-go/automaxprocs/issues/66. This lead to **gomaxecs**.

## Installation

`go get -u github.com/rdforte/gomaxecs`

## Quick Start

```go
import _ "github.com/rdforte/gomaxecs"

func main() {
  // Your application logic here.
}
```

## Intro to GOMAXPROCS

GOMAXPROCS is an env variable and function from the [runtime package](https://pkg.go.dev/runtime@go1.23.1) that limits the number of operating system threads that can execute user-level Go code simultaneously. If GOMAXPROCS is not set then it will default to [runtime.NumCPU](https://pkg.go.dev/runtime@go1.23.1#NumCPU) which is the number of logical CPU cores available by the current process. For example if I decide to run my Go application on my shiny new 8 core Mac Pro, then GOMAXPROCS will default to 8. We are able to configure the number of system threads our Go application can execute by using the runtime.GOMAXPROCS function to override this default.

## What is CFS

CFS was introduced to the Linux kernel in version [2.6.23](https://kernelnewbies.org/Linux_2_6_23) and is the default process scheduler used in Linux. The main purpose behind CFS is to help ensure that each process gets its own fair share of the CPU proportional to its priority. In Docker every container has access to all the hosts resources, within the limits of the kernel scheduler. Though Docker also provides the means to limit these resources through the modifying the containers cgroup on the host machine.

## Performance implications of running Go in Docker

Lets imagine a scenario where we configure our ECS Task to use 8 CPU's and our container to use 4 vCPU's.

```
{
    "containerDefinitions": [
        {
            "cpu": 4096, // Limit container to 4 vCPU's
        }
    ],
    "cpu": "8192", // Task uses 8 CPU's
    "memory": "16384",
    "runtimePlatform": {
        "cpuArchitecture": "X86_64",
        "operatingSystemFamily": "LINUX"
    },
}
```

The ECS Task CPU period is locked into 100ms 

[https://github.com/aws/amazon-ecs-agent/blob/d68e729f73e588982dc2189a1c618c18c47c931b/agent/api/task/task_linux.go#L39](https://github.com/aws/amazon-ecs-agent/blob/d68e729f73e588982dc2189a1c618c18c47c931b/agent/api/task/task_linux.go#L39)

The CPU Period refers to the time period in microseconds, where the kernel will do some calculations to figure out the alloted amount of CPU time to provide each task.
In the above configuration this would be 4 vCPU's multiplied by 100ms giving the task 400ms (4 x 100ms).

If all is well and good with our Go application then we would have go routines scheduled on 4 threads accross 4 cores.

![4 threads](./assets/4-threads.png)

_Threads scheduled on cores 1, 3, 6, 8_

For each 100ms period our Go application consumes the full 400 out of 400ms, therefore 100% of the CPU quota.

Now Go is **NOT** CFS aware https://github.com/golang/go/issues/33803 therefore GOMAXPROCS will default to using all 8 cores of the Task.


![8 threads](./assets/8-threads.png)

Now we have our Go application using all 8 cores resulting in 8 threads executing go routines. After 50ms of execution we reach our CPU quota 50ms * 8 threads giving us 400ms (8 * 50ms).
As a result CFS will throttle our CPU resources, meaning that no more CPU resources will be allocated till the next period. This means our application will be sitting idle doing nothing for
a full 50ms.

If our Go application has an average latency of 50ms this now means a request to our service can up to 150ms which is 300% increase in latency.

## CFS Solution
In Kubernetes this issue is quite easy to solve as we have [uber automaxprocs](https://github.com/uber-go/automaxprocs) to solve this issue. So why not use Uber's automaxprocs then and whats the reason
behind **gomaxecs package**? Well Ubers automaxprocs does not work for ECS https://github.com/uber-go/automaxprocs/issues/66 becuase the cgroup `cpu.cfs_quota_us` is set to -1 🥲. The workaround for this
is to then leverage [ECS Metadata](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint.html) as a means to sourcing the container limts and setting GOMAXPROCS at runtime.

## Contribution
If anyone has any good ideas on how this package can be improved, all contributions are welcome.

## References
[100 Go Mistakes](https://100go.co/?h=kubernetes#not-understanding-the-impacts-of-running-go-in-docker-and-kubernetes-100) was the main source of inspiration for this package. The examples were borrowed from
the book and modified to suit ECS.


<hr>

Released under the [MIT License](LICENSE).


[doc-img]: https://godoc.org/github.com/rdforte/gomaxecs?status.svg
[doc]: https://godoc.org/github.com/rdforte/gomaxecs
[ci-img]: https://github.com/rdforte/gomaxecs/actions/workflows/build.yml/badge.svg
[ci]: https://github.com/rdforte/gomaxecs/actions/workflows/build.yml
[mit-img]: https://img.shields.io/badge/License-MIT-yellow.svg
[mit]: https://github.com/rdforte/gomaxecs/blob/main/LICENSE
