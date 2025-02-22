package repos

import (
	"backend/internal/core/models"
	"backend/internal/integrations"
	"context"
	"database/sql"
	"errors"

	"go.opentelemetry.io/otel/trace"
)

type UserRepo interface {
	CreateUser(ctx context.Context, dto models.UserDTO) (*models.UserDTO, error)
	UpdateUser(ctx context.Context, userId string, dto models.UserUpdateDTO) error
	DeactivateUser(ctx context.Context, userId string) error
	SetUserEmailVerified(ctx context.Context, userId string) error
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

	query := `insert into users (email, secret, full_name) values ($1, $2, $3) returning id;`
	row := u.db.QueryRowContext(ctx, query, dto.Email, dto.Secret, dto.FullName)

	id := ""
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return &models.UserDTO{
		Id:       id,
		Email:    dto.Email,
		Secret:   dto.Secret,
		FullName: dto.FullName,
	}, nil
}

func (u *userRepo) UpdateUser(ctx context.Context, userId string, dto models.UserUpdateDTO) error {
	_, span := u.tracer.Start(ctx, "postgres::UpdateUser")
	defer span.End()

	query := `update users set secret=$1, full_name=$2 where id = $3;`
	_, err := u.db.ExecContext(ctx, query, dto.Secret, dto.FullName, userId)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) DeactivateUser(ctx context.Context, userId string) error {
	_, span := u.tracer.Start(ctx, "postgres::DeactivateUser")
	defer span.End()

	query := `update users set active=false where id = $1;`
	_, err := u.db.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) SetUserEmailVerified(ctx context.Context, userId string) error {
	_, span := u.tracer.Start(ctx, "postgres::SetUserEmailVerified")
	defer span.End()

	query := `update users set email_verified=true where id = $1;`
	_, err := u.db.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) GetUserById(ctx context.Context, id string) (*models.UserDTO, error) {
	_, span := u.tracer.Start(ctx, "postgres::GetUserById")
	defer span.End()

	query := `
	select id, email, secret, full_name, email_verified 
		from users where id = $1 and active;`
	row := u.db.QueryRowContext(ctx, query, id)

	dto := &models.UserDTO{}
	err := row.Scan(&dto.Id, &dto.Email, &dto.Secret, &dto.FullName, &dto.EmailVerified)
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

	query := `select id, email, secret, full_name, email_verified 
		from users where email = $1 and active;`
	row := u.db.QueryRowContext(ctx, query, login)

	dto := &models.UserDTO{}
	err := row.Scan(&dto.Id, &dto.Email, &dto.Secret, &dto.FullName, &dto.EmailVerified)
	if err == nil {
		return dto, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}
