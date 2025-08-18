package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"rest/goAPI/models"

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
	userId, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "invalid jwt token payload",
		})
		return
	}
	typedUserId, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "invalid jwt token payload",
		})
		return
	}
	id, err := strconv.ParseInt(typedUserId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "could not convert id string to int",
		})
		return
	}

	var newEvent models.Event

	err = ctx.ShouldBind(&newEvent)
	newEvent.UserId = id

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
	userId := ctx.GetString("userId")
	intUserId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "could not convert id string to int",
		})
		return
	}

	id := ctx.Param("eventId")

	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid event id",
		})
		return
	}

	var newDetails models.UpdateEvent

	err = ctx.ShouldBind(&newDetails)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid details provided",
		})
		return
	}

	err = models.UpdateEventById(id, newDetails, intUserId)
	if err != nil {
		if err.Error() == "current user does not own this event" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg":   "could not update event",
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "could not update event",
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "event updated successfully",
	})
}

func deleteEvent(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	intUserId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "could not convert id string to int",
		})
		return
	}

	id := ctx.Param("eventId")

	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid event id",
		})
		return
	}

	err = models.DeleteEventById(id, intUserId)
	if err != nil {
		if err.Error() == "current user does not own this event" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg":   "could not update event",
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "could not update event",
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "event deleted successfully",
	})
}
