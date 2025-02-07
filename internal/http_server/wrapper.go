package httpserver

import (
	"backend/pkg/logger"
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type Handler[Input, Output any] func(ctx context.Context, input Input) (Output, error)

type ResponseOk struct {
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}

type ResponseError struct {
	Status string `json:"status"`
	Error  struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	} `json:"error"`
}

func WrapGin[In, Out interface{}](log logger.Logger, handler Handler[In, Out]) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := log.WithContext(c)

		var input In
		if err := c.ShouldBindJSON(&input); err != nil {
			response := ResponseError{
				Status: "error",
				Error: struct {
					Id      string `json:"id"`
					Message string `json:"message"`
				}{
					Id:      "WrongBody",
					Message: err.Error(),
				},
			}

			body, _ := json.Marshal(response)
			c.Data(400, "application/json", body)
			return
		}

		var response interface{}

		output, err := handler(c, input)
		if err != nil {
			log.Error().Err(err).Msg("error in request handler")
			response = ResponseError{
				Status: "error",
				Error: struct {
					Id      string `json:"id"`
					Message string `json:"message"`
				}{
					Id:      "-",
					Message: err.Error(),
				},
			}
		} else {
			response = ResponseOk{
				Status: "success",
				Result: output,
			}
		}

		body, err := json.Marshal(response)
		if err != nil {
			log.Error().Err(err).Msg("marshal response error")
			c.Data(500, "plain/text", []byte(err.Error()))
			return
		}

		c.Data(200, "application/json", body)
	}
}
