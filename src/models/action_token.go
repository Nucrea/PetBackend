package models

type ActionTokenTarget int

const (
	ActionTokenTargetForgotPassword ActionTokenTarget = iota
	ActionTokenTargetLogin2FA
)

type ActionTokenDTO struct {
	Id     string
	UserId string
	Value  string
	Target ActionTokenTarget
}
