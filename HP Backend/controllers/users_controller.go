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
	log.Println(authUrl)
	context.Redirect(http.StatusFound, authUrl)
}

func SpotifyTokenCallBack(context *gin.Context){
	code :=  context.Query("code")
	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Code query parameter is missing"})
		return
	}

	token, err := config.SetSpotifyToken(code)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Token has been stored in DB", "token": token})
}

func TestGetToken(context *gin.Context){
	token, err := config.GetSpotifyTokenObject()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Got Token From db", "token": token})
}