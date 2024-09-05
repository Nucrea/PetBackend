package repos

import (
	"backend/src/core/models"
	"backend/src/integrations"
	"context"
	"database/sql"
)

type ActionTokenRepo interface {
	CreateActionToken(ctx context.Context, dto models.ActionTokenDTO) (*models.ActionTokenDTO, error)
	PopActionToken(ctx context.Context, userId, value string, target models.ActionTokenTarget) (*models.ActionTokenDTO, error)
}

func NewActionTokenRepo(db integrations.SqlDB) ActionTokenRepo {
	return &actionTokenRepo{
		db: db,
	}
}

type actionTokenRepo struct {
	db integrations.SqlDB
}

func (a *actionTokenRepo) CreateActionToken(ctx context.Context, dto models.ActionTokenDTO) (*models.ActionTokenDTO, error) {
	query := `
	insert into 
		action_tokens (user_id, value, target, expiration) 
		values ($1, $2, $3, $4) 
		returning id;`
	row := a.db.QueryRowContext(ctx, query, dto.UserId, dto.Value, dto.Target, dto.Expiration)

	id := ""
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return &models.ActionTokenDTO{
		Id:     id,
		UserId: dto.UserId,
		Value:  dto.Value,
		Target: dto.Target,
	}, nil
}

func (a *actionTokenRepo) PopActionToken(ctx context.Context, userId, value string, target models.ActionTokenTarget) (*models.ActionTokenDTO, error) {
	query := `
	delete 
		from action_tokens 
		where 
			user_id=$1 and value=$2 and target=$3 
			and CURRENT_TIMESTAMP < expiration
		returning id;`
	row := a.db.QueryRowContext(ctx, query, userId, value, target)

	id := ""
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &models.ActionTokenDTO{
		Id:     id,
		UserId: userId,
		Value:  value,
		Target: target,
	}, nil
}
