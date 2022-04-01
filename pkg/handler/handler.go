package handler

import "github.com/gin-gonic/gin"

type Handler struct {
}

func (h *Handler) InitRoutes(ctx *gin.Engine) {
	ctx.GET("/", h.WriteMess)
	ctx.GET("/2", h.WriteMess2)
}
