package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"houseparty.com/models"
	"houseparty.com/services"
	"houseparty.com/websockets"
)

var manager *websockets.Manager

func InitManager(m *websockets.Manager) {
	manager = m
}

func CreateNewRoom(context *gin.Context) {
	var room models.Room

	err := context.ShouldBindJSON(&room)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse data", "error": err.Error()})
		return
	}

	err = services.CreateRoom(&room, context.GetInt64("userId"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create room", "error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "room created", "room": room})
}

func RetieveRooms(context *gin.Context) {

	publicRooms, userRoom, err := services.GetRooms(context.GetInt64("userId"))
	if err != nil {
		log.Println(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrive rooms", "error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "fetched rooms", "public_rooms": publicRooms, "user_room": userRoom})
}

func DeleteRoom(context *gin.Context) {
	roomId := context.Param("id")

	if manager.CountClients(roomId) > 0 {
		errorMessage := fmt.Sprintf("Room has %d active connection(s)", manager.CountClients(roomId))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete room", "error": errorMessage})
		return
	}

	manager.Lock()
	defer manager.Unlock()
	delete(manager.Rooms, roomId)

	room, err := services.DeleteRoomByID(roomId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete room", "error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Room has been deleted.", "room": room})
}

func JoinRoom(context *gin.Context) {
	manager.ServeWs()(context)
}
