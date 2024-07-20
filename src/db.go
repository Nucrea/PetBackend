package src

import (
	"context"
	"database/sql"
	"errors"
)

type DB interface {
	CreateUser(ctx context.Context, dto UserDTO) (*UserDTO, error)
	GetUserById(ctx context.Context, id string) (*UserDTO, error)
	GetUserByLogin(ctx context.Context, login string) (*UserDTO, error)
}

func NewDB(db *sql.DB) DB {
	return &dbImpl{db}
}

type dbImpl struct {
	db *sql.DB
}

func (d *dbImpl) CreateUser(ctx context.Context, dto UserDTO) (*UserDTO, error) {
	query := `insert into users (login, secret, name) values (?, ?, ?) returning id;`
	row := d.db.QueryRowContext(ctx, query, dto.Login, dto.Secret, dto.Name)

	id := ""
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return &UserDTO{
		Id:     id,
		Login:  dto.Login,
		Secret: dto.Secret,
		Name:   dto.Name,
	}, nil
}

func (d *dbImpl) GetUserById(ctx context.Context, id string) (*UserDTO, error) {
	query := `select (id, login, secret, name) from users where id = ?;`
	row := d.db.QueryRowContext(ctx, query, id)

	dto := &UserDTO{}
	err := row.Scan(dto)
	if err == nil {
		return dto, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (d *dbImpl) GetUserByLogin(ctx context.Context, login string) (*UserDTO, error) {
	query := `select (id, login, secret, name) from users where login = ?;`
	row := d.db.QueryRowContext(ctx, query, login)

	dto := &UserDTO{}
	err := row.Scan(dto)
	if err == nil {
		return dto, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}
