package main

import (
	"html/template"
	"strings"

	"gopkg.in/gomail.v2"
)

const MSG_TEXT = `
<html>
	<head>
	</head>
	<body>
		<p>{{.Text}}</p>
		{{if .Link}}
		<a href="{{.Link}}">link</a>
		{{end}}
	</body>
</html>
`

type MailContent struct {
	Text string
	Link string
}

func NewEmailer(conf ConfigSMTP) (*Emailer, error) {
	dialer := gomail.NewDialer(conf.Server, conf.Port, conf.Login, conf.Password)

	closer, err := dialer.Dial()
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	htmlTemplate, err := template.New("email").Parse(MSG_TEXT)
	if err != nil {
		return nil, err
	}

	return &Emailer{
		senderEmail:  conf.Email,
		htmlTemplate: htmlTemplate,
		dialer:       dialer,
	}, nil
}

type Emailer struct {
	senderEmail  string
	htmlTemplate *template.Template
	dialer       *gomail.Dialer
}

func (e *Emailer) SendRestorePassword(email, link string) error {
	return e.sendEmail("Restore your password", email, MailContent{
		Text: "You received this message do request of password change. Use this link to change your password:",
		Link: link,
	})
}

func (e *Emailer) SendVerifyUser(email, link string) error {
	return e.sendEmail("Verify your email", email, MailContent{
		Text: "You received this message due to registration of account. Use this link to verify email:",
		Link: link,
	})
}

func (e *Emailer) SendPasswordChanged(email string) error {
	return e.sendEmail("Password changed", email, MailContent{
		Text: "Your password was successfully changed",
	})
}

func (e *Emailer) sendEmail(subject, to string, content MailContent) error {
	builder := &strings.Builder{}
	if err := e.htmlTemplate.Execute(builder, content); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.senderEmail, "Pet Backend"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", builder.String())

	return e.dialer.DialAndSend(m)
}
