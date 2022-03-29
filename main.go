package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"goServ5/docs"
	"goServ5/pkg/handlers"
	"goServ5/repository/postgres"
)

// @title           CRUD web Server
// @version         1.0
// @description     This is a piece of shit.
// @host      localhost:8888
// @BasePath  /
func main() {
	// initialization postgres db
	postgres.InitDB()
	// programmatically set swagger info
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()

	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/peoples", handlers.GetPeoples)
	router.POST("/peoples/:id", handlers.GetPeoplesById)
	router.POST("/peoples/add", handlers.PostPeoples)
	router.PUT("/peoples", handlers.ModifyPeoples)
	router.DELETE("/peoples/:id", handlers.DeletePeoplesById)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run("localhost:8888")
}
