package repos

import (
	"backend/internal/core/models"
	"backend/internal/integrations"
	"context"
	"database/sql"
)

type ActionTokenRepo interface {
	CreateActionToken(ctx context.Context, dto models.ActionTokenDTO) (*models.ActionTokenDTO, error)
	GetActionToken(ctx context.Context, value string, target models.ActionTokenTarget) (*models.ActionTokenDTO, error)
	DeleteActionToken(ctx context.Context, id string) error
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

func (a *actionTokenRepo) GetActionToken(ctx context.Context, value string, target models.ActionTokenTarget) (*models.ActionTokenDTO, error) {
	dto := &models.ActionTokenDTO{Value: value, Target: target}

	query := `
	select id, user_id from action_tokens 
	where 
		value=$2 and target=$3
		and CURRENT_TIMESTAMP < expiration;`
	row := a.db.QueryRowContext(ctx, query, value, target)

	err := row.Scan(&dto.Id, &dto.UserId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return dto, nil
}

func (a *actionTokenRepo) DeleteActionToken(ctx context.Context, id string) error {
	query := `delete from action_tokens where id=$1;`
	if _, err := a.db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
