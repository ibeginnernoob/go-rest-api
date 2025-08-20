package routes

import (
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
			"msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "events fetched successfully",
		"events": events,
	})
}

func createEvent(ctx *gin.Context) {
	userId, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "no signed in user detected",
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
			"msg": "some error occured in the user's id",
		})
		return
	}

	var newEvent models.Event

	err = ctx.ShouldBind(&newEvent)
	newEvent.UserId = id

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "insufficient event data",
		})
		return
	}
	newEvent.DateTime = time.Now()
	err = newEvent.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"msg":   "event creation successful",
		"event": newEvent,
	})
}

func getEvent(ctx *gin.Context) {
	eventId := ctx.Param("eventId")

	if len(eventId) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "please previde a valid event id",
		})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "event fetched successfully",
		"event": event,
	})
}

func updateEvent(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	intUserId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "some unknown error occured with user id",
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
			"msg": "invalid event details provided",
		})
		return
	}

	err = models.UpdateEventById(id, newDetails, intUserId)
	if err != nil {
		var statusCode int
		if err.Error() == "current user does not own this event" {
			statusCode = http.StatusUnauthorized
		} else {
			statusCode = http.StatusInternalServerError
		}
		ctx.JSON(statusCode, gin.H{
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
			"msg": "some error in the user's id",
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
		var statusCode int
		if err.Error() == "current user does not own this event" {
			statusCode = http.StatusUnauthorized
		} else {
			statusCode = http.StatusInternalServerError
		}
		ctx.JSON(statusCode, gin.H{
			"msg":   "could not update event",
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "event deleted successfully",
	})
}
