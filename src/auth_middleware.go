package src

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware(userService UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("X-Auth")
		if token == "" {
			ctx.AbortWithError(403, fmt.Errorf("authorization required"))
			return
		}

		user, err := userService.ValidateToken(ctx, token)
		if err == ErrUserWrongToken || err == ErrUserNotExists {
			ctx.AbortWithError(403, err)
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
