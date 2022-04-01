package handler

import "github.com/gin-gonic/gin"

//описываем что должны делать хендлеры, без реализации самих методов
//реализация находится в repository.peoples_postgres

func (h *Handler) WriteMess(ctx *gin.Context) {
	ctx.String(200, "1th page")
}
func (h *Handler) WriteMess2(ctx *gin.Context) {
	ctx.String(200, "2th page")
}
