package handlers

import (
	"backend/internal/core/services"
	"backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type A struct {
	Title    string
	Text     string
	Link     string
	LinkText string
}

func NewUserVerifyEmailHandler(log logger.Logger, userService services.UserService) gin.HandlerFunc {
	htmlOk := `
	<html> 
		<head> 
			<title>Verify Email</title>
		</head>
		<body>
			<h1>Email successfuly verified</h1>
		</body> 
	</html>
	`

	htmlNotOk := `
	<html> <head> <title>Verify Email</title> </head> <body>
	<h1>Email was not verified</h1>
	</body> </html>
	`

	return func(c *gin.Context) {
		token, ok := c.GetQuery("token")
		if !ok || token == "" {
			c.Data(400, "text/html", []byte(htmlNotOk))
			return
		}

		err := userService.VerifyEmail(c, token)
		if err != nil {
			log.Error().Err(err).Msg("Error verifying email")
			c.Data(400, "text/html", []byte(htmlNotOk))
			return
		}

		c.Data(200, "text/html", []byte(htmlOk))
	}
}
