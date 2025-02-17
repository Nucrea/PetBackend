package handlers

import (
	"backend/internal/core/services"
	httpserver "backend/internal/http_server"
	"backend/pkg/logger"
	"context"

	"github.com/gin-gonic/gin"
)

type createUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type createUserOutput struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func NewUserCreateHandler(log logger.Logger, userService services.UserService) gin.HandlerFunc {
	return httpserver.WrapGin(log,
		func(ctx context.Context, input createUserInput) (createUserOutput, error) {
			user, err := userService.CreateUser(ctx,
				services.UserCreateParams{
					Email:    input.Email,
					Password: input.Password,
					Name:     input.Name,
				},
			)
			
			if err != nil {
				return createUserOutput{}, err
			}

			return createUserOutput{
				Id:    user.Id,
				Email: user.Email,
				Name:  user.Name,
			}, nil
		},
	)
}
