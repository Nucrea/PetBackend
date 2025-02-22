package models

import "time"

type ActionTokenTarget string

const (
	ActionTokenTargetRestorePassword ActionTokenTarget = "restore"
	ActionTokenTargetVerifyEmail     ActionTokenTarget = "verify"
)

type ActionTokenDTO struct {
	Id         string
	UserId     string
	Value      string
	Target     ActionTokenTarget
	Expiration time.Time
}
