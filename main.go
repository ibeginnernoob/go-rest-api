package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.GET("/", func(context *gin.Context) {

		res := map[string]string{
			"msg": "Hello World",
		}

		fmt.Println(res)

		context.JSON(200, res)
	})
	server.Run(":8080")
}
