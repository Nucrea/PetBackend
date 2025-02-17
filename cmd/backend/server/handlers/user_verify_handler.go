package handlers

import (
	"backend/internal/core/services"
	httpserver "backend/internal/http_server"
	"backend/pkg/logger"
	"context"

	"html/template"

	"github.com/gin-gonic/gin"
)

type HtmlTemplate struct {
	TabTitle string
	Title    string
	Text     string
	Link     string
	LinkText string
}

const htmlTemplate = `
<html> 
	<head>
		<title>{{.TabTitle}}</title>
	</head>
	<body>
		{{if .Title}}
		<h1>{{.Title}}</h1>
		{{end}}

		<h3>{{.Text}}</h3>

		{{if .Link}}
		<a href="{{.Link}}">{{.LinkText}}</a>
		{{end}}
	</body> 
</html>
`

func NewUserVerifyEmailHandler(log logger.Logger, userService services.UserService) gin.HandlerFunc {
	template, err := template.New("verify-email").Parse(htmlTemplate)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing template")
	}

	return func(c *gin.Context) {
		tmp := HtmlTemplate{
			TabTitle: "Verify Email",
			Text:     "Error verifying email",
		}

		token, ok := c.GetQuery("token")
		if !ok || token == "" {
			log.Error().Err(err).Msg("No token in query param")
			template.Execute(c.Writer, tmp)
			c.Status(400)
			return
		}

		err := userService.VerifyEmail(c, token)
		if err != nil {
			log.Error().Err(err).Msg("Error verifying email")
			template.Execute(c.Writer, tmp)
			c.Status(400)
			return
		}

		tmp.Text = "Email successfully verified"
		template.Execute(c.Writer, tmp)
		c.Status(200)
	}
}

type inputSendVerify struct {
	Email string `json:"email" validate:"required,email"`
}

func NewUserSendVerifyEmailHandler(log logger.Logger, userService services.UserService) gin.HandlerFunc {
	return httpserver.WrapGin(log,
		func(ctx context.Context, input inputSendVerify) (interface{}, error) {
			err := userService.SendEmailVerifyEmail(ctx, input.Email)
			return nil, err
		},
	)
}
