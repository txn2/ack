/*
Copyright 2019 txn2

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/txn2/ack"
	"os"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	// gin router
	r := gin.New()

	r.GET("/test", func(c *gin.Context) {
		ak := ack.Gin(c)

		ak.SetPayloadType("Message")
		ak.GinSend("A test message.")
	})

	err := r.Run("127.0.0.1:8080")
	if err != nil {
		print(fmt.Errorf("unable to start server: %s", err.Error()))
		os.Exit(1)
	}
}
