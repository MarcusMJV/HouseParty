package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"houseparty.com/config"
	"houseparty.com/models"
	"houseparty.com/services"
)

func SignUp(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError,  gin.H{"message": "Could not bind user", "error": err.Error()})
		return
	}

	token, err := services.CreateNewUser(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError,  gin.H{"message": "Could not create new user", "error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created", "user": user.ToUserResponse(), "token": token})
}

func Login(context *gin.Context){
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError,  gin.H{"message": "Could not bind user", "error": err.Error()})
		return
	}

	token, err := services.ValidateCredentials(&user)
	if err != nil {
		context.JSON(http.StatusUnauthorized,  gin.H{"message": "Could not validate user", "error": err.Error()})
		return
	}


	context.JSON(http.StatusOK, gin.H{"message": "User has been logged in", "user": user.ToUserResponse(), "token": token})
}

func SpotifyAuthToken(context *gin.Context){
	authUrl, err := config.GenerateSpotifyAuthRequest()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Url Generated", "auth_url": authUrl,})
}

func SpotifyTokenCallBack(context *gin.Context){
	var user models.User
	code := context.Param("code")
	err := user.GetUserById(context.GetInt64("userId")) 
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not get user", "error": err.Error()})
		return
	}

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Code query parameter is missing"})
		return
	}
	token, err := config.SetSpotifyToken(code, user.Id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could get and save token", "error": err.Error()})
		return
	}

	err = user.ActivateSpotify()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not update user", "error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Token has been stored in DB", "token": token})
}

func TestGetToken(context *gin.Context){
	token, err := config.GetSpotifyTokenObject(1)
	log.Println(token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err = config.RefreshToken(token.RefreshToken, 1)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "spotify connected", "token": token})
}