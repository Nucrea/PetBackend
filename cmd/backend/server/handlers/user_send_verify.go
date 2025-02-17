package handlers

import (
	"backend/internal/core/services"
	httpserver "backend/internal/http_server"
	"backend/pkg/logger"
	"context"

	"github.com/gin-gonic/gin"
)

type inputSendVerify struct {
	Email string `json:"email" validate:"required,email"`
}

func NewUserSendVerifyEmailHandler(log logger.Logger, userService services.UserService) gin.HandlerFunc {
	return httpserver.WrapGin(log,
		func(ctx context.Context, input inputSendVerify) (interface{}, error) {
			err := userService.SendEmailVerifyUser(ctx, input.Email)
			return nil, err
		},
	)
}
