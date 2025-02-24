package handlers

import "github.com/gin-gonic/gin"

func New200OkHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(200)
	}
}
