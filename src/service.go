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
	ErrUserBadPassword   = fmt.Errorf("password must contain at least 8 characters")
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
	Db       DB
	Jwt      JwtUtil
	Password PasswordUtil
	Cache    Cache[string, UserDTO]
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

	if err := u.deps.Password.Validate(params.Password); err != nil {
		return nil, ErrUserBadPassword
	}

	secret, err := u.deps.Password.Hash(params.Password)
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

	u.deps.Cache.Set(result.Id, *result, -1)

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

	if !u.deps.Password.Compare(password, user.Secret) {
		return "", ErrUserWrongPassword
	}

	jwt, err := u.deps.Jwt.Create(*user)
	if err != nil {
		return "", err
	}

	u.deps.Cache.Set(user.Id, *user, -1)

	return jwt, nil
}

func (u *userService) getUserById(ctx context.Context, userId string) (*UserDTO, error) {
	if user, ok := u.deps.Cache.Get(userId); ok {
		return &user, nil
	}

	user, err := u.deps.Db.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotExists
	}

	u.deps.Cache.Set(user.Id, *user, -1)

	return user, nil
}

func (u *userService) ValidateToken(ctx context.Context, tokenStr string) (*UserDTO, error) {
	payload, err := u.deps.Jwt.Parse(tokenStr)
	if err != nil {
		return nil, ErrUserWrongToken
	}

	user, err := u.getUserById(ctx, payload.UserId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
