package handlers

import (
	"backend/internal/core/services"
	httpserver "backend/internal/http_server"
	"backend/pkg/logger"
	"context"

	"github.com/gin-gonic/gin"
)

type inputSendRestorePassword struct {
	Email string `json:"email" binding:"required,email"`
}

func NewUserSendRestorePasswordHandler(log logger.Logger, userService services.UserService) gin.HandlerFunc {
	return httpserver.WrapGin(log,
		func(ctx context.Context, input inputSendRestorePassword) (interface{}, error) {
			err := userService.SendEmailForgotPassword(ctx, input.Email)
			return nil, err
		},
	)
}
