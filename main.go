package main

import (
	"rest/goAPI/db"

	"github.com/gin-gonic/gin"

	"rest/goAPI/routes"
)

func main() {
	db.DBInit()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
