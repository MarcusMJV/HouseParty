package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

func TestSpotify(context *gin.Context){
	// apiResponse, err := config.GetSpotifyToken();
	tracks, err := services.SearchSongs("presure paramore")
	if err != nil {
		context.JSON(http.StatusUnauthorized,  gin.H{"message": "error accured", "error": err.Error()})
		return
	}

	songs, err := services.SimplifyTracks(tracks)
	if err != nil {
		context.JSON(http.StatusUnauthorized,  gin.H{"message": "error accured", "error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Spotify Token", "songs": songs})
}