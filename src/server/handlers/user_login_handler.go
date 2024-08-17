package handlers

import (
	"backend/src/core/services"
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

func NewUserLoginHandler(userService services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := loginUserInput{}
		if err := ctx.ShouldBindJSON(&params); err != nil {
			ctx.AbortWithError(400, err)
			return
		}

		token, err := userService.AuthenticateUser(ctx, params.Login, params.Password)
		if err == services.ErrUserNotExists || err == services.ErrUserWrongPassword {
			ctx.AbortWithError(400, err)
			return
		}
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		resultBody, err := json.Marshal(loginUserOutput{
			Token: token,
		})
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		ctx.Data(200, "application/json", resultBody)
	}
}
