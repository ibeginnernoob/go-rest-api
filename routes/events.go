package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"rest/goAPI/models"
	"rest/goAPI/utils"

	"github.com/gin-gonic/gin"
)

func getEvents(ctx *gin.Context) {
	events, err := models.GetEvents()
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
	token := ctx.Request.Header.Get("Authorization")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "no jwt sent",
		})
		return
	}

	payload, err := utils.ValidateToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "invalid jwt, pls sign in",
		})
		return
	}

	var newEvent models.Event
	userId, err := strconv.ParseInt(payload.Id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "could not convert id string to int",
		})
		return
	}

	newEvent.UserId = userId

	err = ctx.ShouldBind(&newEvent)

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

	event, err := models.GetEventByID(eventId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "could not fetch event",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "event fetch successfully",
		"event": event,
	})
}

func updateEvent(ctx *gin.Context) {
	id := ctx.Param("eventId")

	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid event id",
		})
		return
	}

	var newDetails models.UpdateEvent

	err := ctx.ShouldBind(&newDetails)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid details provided",
		})
		return
	}

	err = models.UpdateEventById(id, newDetails)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "could not update event",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "event updated successfully",
	})
}

func deleteEvent(ctx *gin.Context) {
	id := ctx.Param("eventId")

	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid event id",
		})
		return
	}

	err := models.DeleteEventById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "could not delete event",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "event deleted successfully",
	})
}
