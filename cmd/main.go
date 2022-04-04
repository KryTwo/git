package main

import (
	_ "github.com/lib/pq"
	"goServ5/docs"
	"goServ5/pkg/handlers"
	"goServ5/repository/postgres"
)

// @title           CRUD web Server
// @version         1.0
// @host      localhost:8888
// @BasePath  /
func main() {
	// initialization postgres db
	postgres.InitDB()
	// programmatically set swagger info
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	// initializing routes
	handlers.InitRoutes()
}
