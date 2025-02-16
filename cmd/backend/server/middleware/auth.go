package middleware

import (
	"backend/internal/core/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware(userService services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("X-Auth")
		if token == "" {
			ctx.AbortWithError(403, fmt.Errorf("authorization required"))
			return
		}

		user, err := userService.ValidateAuthToken(ctx, token)
		if err == services.ErrUserWrongToken || err == services.ErrUserNotExists {
			ctx.AbortWithError(403, err)
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
