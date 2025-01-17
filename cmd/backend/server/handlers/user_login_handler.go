package handlers

import (
	"backend/internal/core/services"
	"backend/pkg/logger"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type loginUserInput struct {
	Login    string `json:"email"`
	Password string `json:"password"`
}

type loginUserOutput struct {
	Token string `json:"token"`
}

func NewUserLoginHandler(logger logger.Logger, userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxLogger := logger.WithContext(c).WithPrefix("NewUserLoginHandler")

		params := loginUserInput{}
		if err := c.ShouldBindJSON(&params); err != nil {
			ctxLogger.Error().Err(err).Msg("bad input body model")
			c.AbortWithError(400, err)
			return
		}

		token, err := userService.AuthenticateUser(c, params.Login, params.Password)
		if err == services.ErrUserNotExists {
			ctxLogger.Error().Err(err).Msg("user does not exist")
			c.AbortWithError(400, err)
			return
		}
		if err == services.ErrUserWrongPassword {
			ctxLogger.Error().Err(err).Msg("wrong password")
			c.AbortWithError(400, err)
			return
		}
		if err != nil {
			ctxLogger.Error().Err(err).Msg("AuthenticateUser internal error")
			c.AbortWithError(500, err)
			return
		}

		resultBody, err := json.Marshal(loginUserOutput{
			Token: token,
		})
		if err != nil {
			ctxLogger.Error().Err(err).Msg("marshal json internal error")
			c.AbortWithError(500, err)
			return
		}

		c.Data(200, "application/json", resultBody)
	}
}
