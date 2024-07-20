package src

import (
	"context"
	"fmt"
)

var (
	ErrUserNotExists     = fmt.Errorf("no such user")
	ErrUserExists        = fmt.Errorf("user with this login already exists")
	ErrUserWrongPassword = fmt.Errorf("wrong password")
	ErrUserWrongToken    = fmt.Errorf("bad user token")
	// ErrUserInternal = fmt.Errorf("unexpected error. contact tech support")
)

type UserService interface {
	CreateUser(ctx context.Context, params UserCreateParams) (*UserDTO, error)
	AuthenticateUser(ctx context.Context, login, password string) (string, error)
	ValidateToken(ctx context.Context, tokenStr string) (*UserDTO, error)
}

func NewUserService(deps UserServiceDeps) UserService {
	return &userService{deps}
}

type UserServiceDeps struct {
	Db     DB
	Jwt    JwtUtil
	Bcrypt BCryptUtil
}

type userService struct {
	deps UserServiceDeps
}

type UserCreateParams struct {
	Login    string
	Password string
	Name     string
}

func (u *userService) CreateUser(ctx context.Context, params UserCreateParams) (*UserDTO, error) {
	exisitngUser, err := u.deps.Db.GetUserByLogin(ctx, params.Login)
	if err != nil {
		return nil, err
	}
	if exisitngUser != nil {
		return nil, ErrUserExists
	}

	secret, err := u.deps.Bcrypt.HashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	user := UserDTO{
		Login:  params.Login,
		Secret: string(secret),
		Name:   params.Name,
	}

	result, err := u.deps.Db.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *userService) AuthenticateUser(ctx context.Context, login, password string) (string, error) {
	user, err := u.deps.Db.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUserNotExists
	}

	if !u.deps.Bcrypt.IsPasswordsEqual(password, user.Secret) {
		return "", ErrUserWrongPassword
	}

	jwt, err := u.deps.Jwt.Create(*user)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (u *userService) ValidateToken(ctx context.Context, tokenStr string) (*UserDTO, error) {
	payload, err := u.deps.Jwt.Parse(tokenStr)
	if err != nil {
		return nil, ErrUserWrongToken
	}

	user, err := u.deps.Db.GetUserById(ctx, payload.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotExists
	}

	return user, nil
}
