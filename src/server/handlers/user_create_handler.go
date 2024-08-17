package handlers

import (
	"backend/src/core/services"
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

func NewUserCreateHandler(userService services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := createUserInput{}
		if err := ctx.ShouldBindJSON(&params); err != nil {
			ctx.AbortWithError(400, err)
			return
		}

		dto, err := userService.CreateUser(ctx, services.UserCreateParams{
			Email:    params.Email,
			Password: params.Password,
			Name:     params.Name,
		})
		if err == services.ErrUserExists || err == services.ErrUserBadPassword {
			ctx.Data(400, "plain/text", []byte(err.Error()))
			return
		}
		if err != nil {
			ctx.Data(500, "plain/text", []byte(err.Error()))
			return
		}

		resultBody, err := json.Marshal(createUserOutput{
			Id:    dto.Id,
			Email: dto.Email,
			Name:  dto.Name,
		})
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		ctx.Data(200, "application/json", resultBody)
	}
}
