package main

import (
	"github.com/gin-gonic/gin"
	"houseparty.com/config"
	"houseparty.com/controllers"
	"houseparty.com/middleware"
	"houseparty.com/routes"
	"houseparty.com/storage"
	"houseparty.com/websockets"
)

func main() {
	storage.InitDB()
	config.LoadEnv()


	manager := websockets.NewManager()
	controllers.InitManager(manager)
	

	server := gin.Default()
	server.Use(middleware.Cors())
	routes.RegisterRoutes(server)
	server.Run(":8080")
}