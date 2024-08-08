package repos

import (
	"strings"

	"gopkg.in/gomail.v2"
)

const MSG_TEXT = `
<html>
	<head>
	</head>
	<body>
		<p>This message was sent because you forgot a password</p>
		<p>To change a password, use <a href="{{Link}}"/>this</a> link</p>
	</body>
</html>
`

type EmailRepo interface {
	SendEmailForgotPassword(email, token string)
}

func NewEmailRepo() EmailRepo {
	return &emailRepo{}
}

type emailRepo struct {
	// mail *gomail.Dialer
}

func (e *emailRepo) SendEmailForgotPassword(email, token string) {
	link := "https://nucrea.ru?token=" + token
	msgText := strings.ReplaceAll(MSG_TEXT, "{{Link}}", link)

	m := gomail.NewMessage()
	m.SetHeader("From", "email")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", msgText)

	d := gomail.NewDialer("smtp.yandex.ru", 587, "login", "password")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
