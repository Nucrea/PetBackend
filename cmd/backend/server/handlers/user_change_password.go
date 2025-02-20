package handlers

import (
	"backend/cmd/backend/server/middleware"
	"backend/internal/core/services"
	httpserver "backend/internal/http_server"
	"backend/pkg/logger"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type inputChangePassword struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func NewUserChangePasswordHandler(log logger.Logger, userService services.UserService) gin.HandlerFunc {
	return httpserver.WrapGin(log,
		func(ctx context.Context, input inputChangePassword) (interface{}, error) {
			ginCtx, ok := ctx.(*gin.Context)
			if !ok {
				return nil, fmt.Errorf("can not cast context")
			}
			user := middleware.GetUserFromRequest(ginCtx)

			err := userService.ChangePassword(ctx, user.Id, input.OldPassword, input.NewPassword)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
}
