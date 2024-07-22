package handlers

import "github.com/gin-gonic/gin"

func NewDummyHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(200)
	}
}
