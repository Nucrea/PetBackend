package repos

import (
	"backend/src/core/models"
	"backend/src/integrations"
	"context"
	"database/sql"
	"errors"

	"go.opentelemetry.io/otel/trace"
)

// type userDAO struct {
// 	Id     string `json:"id"`
// 	Login  string `json:"login"`
// 	Secret string `json:"secret"`
// 	Name   string `json:"name"`
// }

type UserRepo interface {
	CreateUser(ctx context.Context, dto models.UserDTO) (*models.UserDTO, error)
	UpdateUser(ctx context.Context, userId string, dto models.UserUpdateDTO) error
	GetUserById(ctx context.Context, id string) (*models.UserDTO, error)
	GetUserByEmail(ctx context.Context, login string) (*models.UserDTO, error)
}

func NewUserRepo(db integrations.SqlDB, tracer trace.Tracer) UserRepo {
	return &userRepo{db, tracer}
}

type userRepo struct {
	db     integrations.SqlDB
	tracer trace.Tracer
}

func (u *userRepo) CreateUser(ctx context.Context, dto models.UserDTO) (*models.UserDTO, error) {
	_, span := u.tracer.Start(ctx, "postgres::CreateUser")
	defer span.End()

	query := `insert into users (email, secret, name) values ($1, $2, $3) returning id;`
	row := u.db.QueryRowContext(ctx, query, dto.Email, dto.Secret, dto.Name)

	id := ""
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return &models.UserDTO{
		Id:     id,
		Email:  dto.Email,
		Secret: dto.Secret,
		Name:   dto.Name,
	}, nil
}

func (u *userRepo) UpdateUser(ctx context.Context, userId string, dto models.UserUpdateDTO) error {
	_, span := u.tracer.Start(ctx, "postgres::UpdateUser")
	defer span.End()

	query := `update users set secret=$1, name=$2 where id = $3;`
	_, err := u.db.ExecContext(ctx, query, dto.Secret, dto.Name, userId)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) GetUserById(ctx context.Context, id string) (*models.UserDTO, error) {
	_, span := u.tracer.Start(ctx, "postgres::GetUserById")
	defer span.End()

	query := `select id, email, secret, name from users where id = $1;`
	row := u.db.QueryRowContext(ctx, query, id)

	dto := &models.UserDTO{}
	err := row.Scan(&dto.Id, &dto.Email, &dto.Secret, &dto.Name)
	if err == nil {
		return dto, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (u *userRepo) GetUserByEmail(ctx context.Context, login string) (*models.UserDTO, error) {
	_, span := u.tracer.Start(ctx, "postgres::GetUserByEmail")
	defer span.End()

	query := `select id, email, secret, name from users where email = $1;`
	row := u.db.QueryRowContext(ctx, query, login)

	dto := &models.UserDTO{}
	err := row.Scan(&dto.Id, &dto.Email, &dto.Secret, &dto.Name)
	if err == nil {
		return dto, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}
