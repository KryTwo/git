package server

import (
	"github.com/gin-gonic/gin"
	"main/pkg/handler"
)

type Server struct {
	*gin.Engine
}

type Runner interface {
	Run()
}

func (s *Server) Run() {

	var h *handler.Handler

	router := gin.New()
	h.InitRoutes(router)

	router.Run(":8888")
}
