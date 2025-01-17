package handlers

import (
	"backend/internal/core/services"
	"backend/pkg/logger"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type createUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type createUserOutput struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func NewUserCreateHandler(logger logger.Logger, userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxLogger := logger.WithContext(c)

		params := createUserInput{}
		if err := c.ShouldBindJSON(&params); err != nil {
			ctxLogger.Error().Err(err).Msg("bad input body model")
			c.Data(400, "plain/text", []byte(err.Error()))
			return
		}

		dto, err := userService.CreateUser(
			c,
			services.UserCreateParams{
				Email:    params.Email,
				Password: params.Password,
				Name:     params.Name,
			},
		)
		if err == services.ErrUserExists {
			ctxLogger.Error().Err(err).Msg("user already exists")
			c.Data(400, "plain/text", []byte(err.Error()))
			return
		}
		if err == services.ErrUserBadPassword {
			ctxLogger.Error().Err(err).Msg("password does not satisfy requirements")
			c.Data(400, "plain/text", []byte(err.Error()))
			return
		}
		if err != nil {
			ctxLogger.Error().Err(err).Msg("unexpected create user error")
			c.Data(500, "plain/text", []byte(err.Error()))
			return
		}

		resultBody, err := json.Marshal(
			createUserOutput{
				Id:    dto.Id,
				Email: dto.Email,
				Name:  dto.Name,
			},
		)
		if err != nil {
			ctxLogger.Error().Err(err).Msg("marshal user model error")
			c.Data(500, "plain/text", []byte(err.Error()))
			return
		}

		c.Data(200, "application/json", resultBody)
	}
}
