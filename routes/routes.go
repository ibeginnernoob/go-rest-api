package routes

import (
	"net/http"
	"rest/goAPI/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "Hello World!",
		})
	})
	server.POST("/signup", signup)
	server.POST("/signin", signin)

	server.GET("/events", getEvents)
	server.GET("/events/:eventId", getEvent)
	authorized := server.Group("/")
	authorized.Use(middlewares.IsAuth)
	{
		authorized.POST("/events", createEvent)
		authorized.PUT("/events/:eventId", updateEvent)
		authorized.DELETE("/events/:eventId", deleteEvent)
		authorized.POST("/register", createRegistraion)
		authorized.DELETE("/register/:registrationId", deleteRegistration)
	}

	server.GET("/users", getUsers)
	server.DELETE("/users", deleteUsers)
}
