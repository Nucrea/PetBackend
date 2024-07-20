package src

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type createUserInput struct {
	Login    string
	Password string
	Name     string
}

type createUserOutput struct {
	Id    string `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

func NewUserCreateHandler(userService UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := createUserInput{}
		if err := ctx.ShouldBindJSON(&params); err != nil {
			ctx.AbortWithError(400, err)
			return
		}

		dto, err := userService.CreateUser(ctx, UserCreateParams{
			Login:    params.Login,
			Password: params.Password,
			Name:     params.Name,
		})
		if err == ErrUserExists || err == ErrUserBadPassword {
			ctx.Data(400, "plain/text", []byte(err.Error()))
			return
		}
		if err != nil {
			ctx.Data(500, "plain/text", []byte(err.Error()))
			return
		}

		resultBody, err := json.Marshal(createUserOutput{
			Id:    dto.Id,
			Login: dto.Login,
			Name:  dto.Name,
		})
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		ctx.Data(200, "application/json", resultBody)
	}
}
