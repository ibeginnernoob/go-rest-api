package main

import (
	"fmt"
	"net/http"
	"time"

	"rest/goAPI/db"
	"rest/goAPI/models"

	"github.com/gin-gonic/gin"
)

func main() {
	db.DBInit()
	server := gin.Default()

	server.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "Hello World!",
		})
	})
	server.GET("/events", getEvents)
	server.GET("/event/:eventId", getEvent)
	server.POST("/events", createEvent)

	server.Run(":8080")
}

func getEvents(ctx *gin.Context) {

	events, err := models.GetEvents()
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "could not fetch events",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "done",
		"events": events,
	})
}

func createEvent(ctx *gin.Context) {
	var newEvent models.Event

	err := ctx.ShouldBind(&newEvent)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "Insufficient event data provided.",
		})
		return
	}
	newEvent.DateTime = time.Now()
	err = newEvent.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"msg":   "Event creation successful.",
		"event": newEvent,
	})
}

func getEvent(ctx *gin.Context) {
	eventId := ctx.Param("eventId")
	fmt.Println("new request with id", eventId)

	if len(eventId) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "pls enter a valid event id in the req params",
		})
		return
	}

	event, err := models.GetEvent(eventId)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "could not fetch event",
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"msg":   "event fetch successfully",
		"event": event,
	})
}
