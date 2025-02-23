package routes

import (
	"github.com/gin-gonic/gin"
	"houseparty.com/controllers"
	"houseparty.com/middleware"
)

func RegisterRoutes(server *gin.Engine){
	authenticated := server.Group("/")
	authenticated.Use(middleware.Authenticate)

	authenticated.POST("/room/create",controllers.CreateNewRoom )
	authenticated.GET("/rooms",controllers.RetieveRooms )
	authenticated.GET("/join/room/:id",controllers.JoinRoom )
	authenticated.DELETE("/room/delete/:id", controllers.DeleteRoom)
	authenticated.GET("/auth/token", controllers.SpotifyAuthToken)
	authenticated.POST("/spotify/token/callback/:code", controllers.SpotifyTokenCallBack)
	
	server.POST("/signup", controllers.SignUp)
	server.POST("/login", controllers.Login)
	server.GET("/get/token", controllers.TestGetToken)

}