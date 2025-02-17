package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"houseparty.com/utils"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		token = context.Query("token")
	}else{
		token = strings.Replace(token, "Bearer ", "", 1)
		token = strings.TrimSpace(token)
	}
	
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token Not Found"})
		return
	}

	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token not Valid",  "error": err})
		return
	}

	context.Set("userId", userId)
	context.Next()
}