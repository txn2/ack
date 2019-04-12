# Ack

TXN2 common API acknowledgement (ACK) struct and wrapped gin-gonic Go library.

Ack server wraps gin-gonic and contains flag and environment variable configuration for
providing health and metrics endpoints:
- Basic Auth protected `/healthz` for liveness probes.
- Prometheus metrics on port **2112** at `/metrics`.

## Example Implementation

```go
package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/txn2/ack"
)

func main() {

	test := flag.Bool("test", true, "A test flag")

	server := ack.NewServer()

	if *test {
		server.Router.GET("/test", func(c *gin.Context) {
			ak := ack.Gin(c)
			ak.SetPayloadType("Message")
			ak.GinSend("A test message.")
		})
	}

	server.Run()
}
```

## Run Example

```bash
AGENT=test SERVICE_ENV=dev SERVICE_NS=example go run ./example/server.go --test
```

## Configuration


| Flag          | Environment Variable | Description                                            |
|:--------------|:---------------------|:-------------------------------------------------------|
| -help         |                      | Display help                                           |
| -debug        | DEBUG                | Debug logging mode (default false)                     |
| -ip           | IP                   | Server IP address to bind to. (default "127.0.0.1")    |
| -port         | PORT                 | Server port. (default "8080")                          |
| -healthz      | HEALTHZ              | Enable or disable /healthz (default true)              |
| -healthzUser  | HEALTHZ_USER         | /healthz basic auth username (default "healthz")       |
| -healthzPass  | HEALTHZ_PASS         | /healthz basic auth password (default "healthz")       |
| -logout       | LOGOUT               | log output stdout \|  (default "stdout")               |
| -metric       | METRICS              | Enable or Disable metrics (default true)               |
| -metricsIP    | METRICS_IP           | Falls back to same IP as server. (default "127.0.0.1") |
| -metricsPort  | METRICS_PORT         | Metrics port. (default "2112")                         |
| -readTimeout  | READ_TIMEOUT         | HTTP read timeout in seconds (default 10)              |
| -writeTimeout | WRITE_TIMEOUT        | HTTP write timeout (default 10)                                                       |
|               | AGENT                | Populates the agent key of Ack.                        |
|               | SERVICE_ENV          | Populates the srv_env key of Ack.                      |
|               | SERVICE_NS           | Populates the srv_ns key of Ack.                       |
