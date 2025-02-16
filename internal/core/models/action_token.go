package models

import "time"

type ActionTokenTarget int

const (
	_ ActionTokenTarget = iota
	ActionTokenTargetForgotPassword
	ActionTokenTargetLogin2FA
	ActionTokenVerifyEmail
)

type ActionTokenDTO struct {
	Id         string
	UserId     string
	Value      string
	Target     ActionTokenTarget
	Expiration time.Time
}
