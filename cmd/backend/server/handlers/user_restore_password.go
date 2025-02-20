package handlers

import (
	"backend/internal/core/services"
	httpserver "backend/internal/http_server"
	"backend/pkg/logger"
	"context"

	"github.com/gin-gonic/gin"
)

type inputRestorePassword struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"password" binding:"required"`
}

func NewUserRestorePasswordHandler(log logger.Logger, userService services.UserService) gin.HandlerFunc {
	return httpserver.WrapGin(log,
		func(ctx context.Context, input inputRestorePassword) (interface{}, error) {
			err := userService.ChangePasswordWithToken(ctx, input.Token, input.NewPassword)
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
	)
}
