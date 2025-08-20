package routes

import (
	"net/http"
	"rest/goAPI/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createRegistraion(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	intUserId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "some error in the user's id",
		})
		return
	}

	var registration models.Registration
	err = ctx.Bind(&registration)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid request body",
		})
		return
	}

	_, err = models.GetEventByID(strconv.Itoa(int(registration.EventId)))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid event id",
		})
		return
	}

	registration.UserId = intUserId

	id, err := registration.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"msg": "registration successful",
		"id":  id,
	})
}

func deleteRegistration(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	intUserId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "some error in the user's id",
		})
		return
	}

	var registrationId int64
	registrationId, err = strconv.ParseInt(ctx.Param("registrationId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "some error in the registration id",
		})
		return
	}

	err = models.DeleteRegistration(registrationId, intUserId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "registration could not be deleted, please try again later",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "registration deleted successfully",
	})
}
