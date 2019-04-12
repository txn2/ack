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
