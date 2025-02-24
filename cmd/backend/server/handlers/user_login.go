package handlers

import (
	"backend/internal/core/services"
	httpserver "backend/internal/http_server"
	"backend/pkg/logger"
	"context"

	"github.com/gin-gonic/gin"
)

type loginUserInput struct {
	Login    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginUserOutput struct {
	Token string `json:"token"`
}

func NewUserLoginHandler(log logger.Logger, userService services.UserService) gin.HandlerFunc {
	return httpserver.WrapGin(log,
		func(ctx context.Context, input loginUserInput) (loginUserOutput, error) {
			token, err := userService.AuthenticateUser(ctx, input.Login, input.Password)
			if err != nil {
				return loginUserOutput{}, err
			}

			return loginUserOutput{
				Token: token,
			}, nil
		},
	)
}
