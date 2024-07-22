package services

import (
	"backend/src/models"
	"backend/src/repo"
	"backend/src/utils"
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
	CreateUser(ctx context.Context, params UserCreateParams) (*models.UserDTO, error)
	AuthenticateUser(ctx context.Context, login, password string) (string, error)
	ValidateToken(ctx context.Context, tokenStr string) (*models.UserDTO, error)
}

func NewUserService(deps UserServiceDeps) UserService {
	return &userService{deps}
}

type UserServiceDeps struct {
	Jwt       utils.JwtUtil
	Password  utils.PasswordUtil
	UserRepo  repo.UserRepo
	UserCache repo.Cache[string, models.UserDTO]
}

type userService struct {
	deps UserServiceDeps
}

type UserCreateParams struct {
	Login    string
	Password string
	Name     string
}

func (u *userService) CreateUser(ctx context.Context, params UserCreateParams) (*models.UserDTO, error) {
	exisitngUser, err := u.deps.UserRepo.GetUserByLogin(ctx, params.Login)
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

	user := models.UserDTO{
		Login:  params.Login,
		Secret: string(secret),
		Name:   params.Name,
	}

	result, err := u.deps.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	u.deps.UserCache.Set(result.Id, *result, -1)

	return result, nil
}

func (u *userService) AuthenticateUser(ctx context.Context, login, password string) (string, error) {
	user, err := u.deps.UserRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUserNotExists
	}

	if !u.deps.Password.Compare(password, user.Secret) {
		return "", ErrUserWrongPassword
	}

	payload := utils.JwtPayload{UserId: user.Id}
	jwt, err := u.deps.Jwt.Create(payload)
	if err != nil {
		return "", err
	}

	u.deps.UserCache.Set(user.Id, *user, -1)

	return jwt, nil
}

func (u *userService) getUserById(ctx context.Context, userId string) (*models.UserDTO, error) {
	if user, ok := u.deps.UserCache.Get(userId); ok {
		return &user, nil
	}

	user, err := u.deps.UserRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotExists
	}

	u.deps.UserCache.Set(user.Id, *user, -1)

	return user, nil
}

func (u *userService) ValidateToken(ctx context.Context, tokenStr string) (*models.UserDTO, error) {
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
