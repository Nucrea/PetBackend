package repo

import (
	"backend/src/models"
	"context"
	"database/sql"
	"errors"
)

// type userDAO struct {
// 	Id     string `json:"id"`
// 	Login  string `json:"login"`
// 	Secret string `json:"secret"`
// 	Name   string `json:"name"`
// }

type UserRepo interface {
	CreateUser(ctx context.Context, dto models.UserDTO) (*models.UserDTO, error)
	GetUserById(ctx context.Context, id string) (*models.UserDTO, error)
	GetUserByLogin(ctx context.Context, login string) (*models.UserDTO, error)
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{db}
}

type userRepo struct {
	db *sql.DB
}

func (u *userRepo) CreateUser(ctx context.Context, dto models.UserDTO) (*models.UserDTO, error) {
	query := `insert into users (login, secret, name) values ($1, $2, $3) returning id;`
	row := u.db.QueryRowContext(ctx, query, dto.Login, dto.Secret, dto.Name)

	id := ""
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return &models.UserDTO{
		Id:     id,
		Login:  dto.Login,
		Secret: dto.Secret,
		Name:   dto.Name,
	}, nil
}

func (u *userRepo) GetUserById(ctx context.Context, id string) (*models.UserDTO, error) {
	query := `select id, login, secret, name from users where id = $1;`
	row := u.db.QueryRowContext(ctx, query, id)

	dto := &models.UserDTO{}
	err := row.Scan(&dto.Id, &dto.Login, &dto.Secret, &dto.Name)
	if err == nil {
		return dto, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (u *userRepo) GetUserByLogin(ctx context.Context, login string) (*models.UserDTO, error) {
	query := `select id, login, secret, name from users where login = $1;`
	row := u.db.QueryRowContext(ctx, query, login)

	dto := &models.UserDTO{}
	err := row.Scan(&dto.Id, &dto.Login, &dto.Secret, &dto.Name)
	if err == nil {
		return dto, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}
