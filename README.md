# gomaxecs [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci]

Due to Go not being CFS aware 🥲 https://github.com/golang/go/issues/33803 and because [uber automaxprocs](https://github.com/uber-go/automaxprocs) is unable to set GOMAXPROCS for ECS https://github.com/uber-go/automaxprocs/issues/66. This lead to **gomaxecs** a package that auto sets GOMAXPROCS to the ECS Container Limit.

## Installation

`go get -u github.com/rdforte/gomaxecs`

## Quick Start

```go
import _ "github.com/rdforte/gomaxecs"

func main() {
  // Your application logic here.
}
```

## GOMAXPROCS Performance issues in ECS

GOMAXPROCS is an env variable and function from the [runtime package](https://pkg.go.dev/runtime@go1.23.1) that limits the number of operating system threads that can execute user-level Go code simultaneously. If GOMAXPROCS is not set then it will default to [runtime.NumCPU](https://pkg.go.dev/runtime@go1.23.1#NumCPU) which is the number of logical CPU cores available by the current process. For example if I decide to run my Go application on my shiny new 8 core Mac Pro, then GOMAXPROCS will default to 8. We are able to configure the number of system threads our Go application can execute by using the runtime.GOMAXPROCS function to override this default.

To understand the issue GOMAXPROCS can impose on our Go applications when we run our Go applications in Docker we must also look into a concept referred to as Completely Fair Scheduler or CFS for short.
CFS was introduced to the Linux kernel in version [2.6.23](https://kernelnewbies.org/Linux_2_6_23) and is the default process scheduler used in Linux. The main purpose behind CFS is to help ensure that each process gets its own fair share of the CPU proportional to its priority.

Docker relies on the use of CFS to manage how much CPU a container.
TODO finish


<hr>

Released under the [MIT License](LICENSE).


[doc-img]: https://godoc.org/github.com/rdforte/gomaxecs?status.svg
[doc]: https://godoc.org/github.com/rdforte/gomaxecs
[ci-img]: https://github.com/rdforte/gomaxecs/actions/workflows/build.yml/badge.svg
[ci]: https://github.com/rdforte/gomaxecs/actions/workflows/build.yml
