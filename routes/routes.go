package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "Hello World!",
		})
	})
	server.GET("/events", getEvents)
	server.GET("/events/:eventId", getEvent)
	server.POST("/events", createEvent)
	server.PUT("/events/:eventId", updateEvent)
	server.DELETE("/events/:eventId", deleteEvent)
}
