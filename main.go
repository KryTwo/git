package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"goServ5/pkg/handlers"
	"goServ5/repository/postgres"
)

func main() {
	postgres.InitDB()

	router := gin.Default()
	router.GET("/peoples/", handlers.GetPeoples)
	router.GET("/peoples/:id", handlers.GetPeoplesById)
	router.POST("/peoples/", handlers.PostPeoples)
	router.PUT("/peoples/", handlers.ModifyPeoples)
	router.DELETE("/peoples/:id", handlers.DeletePeoplesById)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run("localhost:8989")
}
