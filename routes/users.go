package routes

import (
	"net/http"
	"rest/goAPI/models"
	"rest/goAPI/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SigninDetails struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func getUsers(ctx *gin.Context) {
	users, err := models.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "could not fetch users",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "users fetched successfully",
		"users": users,
	})
}

func deleteUsers(ctx *gin.Context) {
	err := models.DeleteUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "could not delete users",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "users deleted successfully",
	})
}

func signup(ctx *gin.Context) {
	var newUser models.User

	err := ctx.Bind(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid user details",
		})
		return
	}

	userId, err := newUser.SaveUser()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "could not create new user",
		})
		return
	}

	token, err := utils.CreateToken(strconv.FormatInt(userId, 10))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "could not create jwt",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"msg":   "user sign up success",
		"token": token,
	})
}

func signin(ctx *gin.Context) {
	var details SigninDetails
	err := ctx.Bind(&details)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid fields",
		})
		return
	}

	user, err := models.GetUserByEmail(details.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "invalid email",
		})
		return
	}

	if !utils.CheckPasswordHash(details.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "invalid password",
		})
		return
	}

	token, err := utils.CreateToken(strconv.FormatInt(user.Id, 10))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "could not create jwt",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "user sign in success",
		"token": token,
	})
}
