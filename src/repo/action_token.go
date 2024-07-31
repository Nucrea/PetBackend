package repo

import "backend/src/models"

type ActionTokenRepo interface {
	CreateActionToken(actionToken models.ActionTokenDTO) (*models.ActionTokenDTO, error)
	FindActionToken(userId, val string, target models.ActionTokenTarget) (*models.ActionTokenDTO, error)
	DeleteActionToken(id string) error
}
