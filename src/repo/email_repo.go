package repo

type EmailRepo interface {
	SendEmailForgotPassword(email, token string)
}
