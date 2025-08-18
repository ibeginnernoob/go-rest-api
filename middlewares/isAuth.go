package middlewares

import (
	"net/http"
	"rest/goAPI/utils"

	"github.com/gin-gonic/gin"
)

func IsAuth(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("Authorization")
	payload, err := utils.ValidateToken(tokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	userId := payload.Id
	if userId == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	ctx.Set("userId", userId)
	ctx.Next()
}
