package services

import (
	"backend/internal/core/models"
	"backend/internal/core/repos"
	"backend/internal/core/utils"
	"backend/pkg/cache"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotExists     = fmt.Errorf("no such user")
	ErrUserExists        = fmt.Errorf("user with this login already exists")
	ErrUserWrongPassword = fmt.Errorf("wrong password")
	ErrUserWrongToken    = fmt.Errorf("bad user token")
	ErrUserBadPassword   = fmt.Errorf("password must contain at least 8 characters")
	// ErrUserInternal = fmt.Errorf("unexpected error. contact tech support")
)

const (
	userCacheTtl = time.Hour
)

type UserService interface {
	CreateUser(ctx context.Context, params UserCreateParams) (*models.UserDTO, error)
	AuthenticateUser(ctx context.Context, login, password string) (string, error)
	ValidateToken(ctx context.Context, tokenStr string) (*models.UserDTO, error)

	ForgotPassword(ctx context.Context, userId string) error
	ChangePassword(ctx context.Context, userId, oldPassword, newPassword string) error
	ChangePasswordWithToken(ctx context.Context, userId, actionToken, newPassword string) error
}

func NewUserService(deps UserServiceDeps) UserService {
	return &userService{deps}
}

type UserServiceDeps struct {
	Jwt             utils.JwtUtil
	Password        utils.PasswordUtil
	UserRepo        repos.UserRepo
	UserCache       cache.Cache[string, models.UserDTO]
	JwtCache        cache.Cache[string, string]
	EventRepo       repos.EventRepo
	ActionTokenRepo repos.ActionTokenRepo
}

type userService struct {
	deps UserServiceDeps
}

type UserCreateParams struct {
	Email    string
	Password string
	Name     string
}

func (u *userService) CreateUser(ctx context.Context, params UserCreateParams) (*models.UserDTO, error) {
	exisitngUser, err := u.deps.UserRepo.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return nil, err
	}
	if exisitngUser != nil {
		return nil, ErrUserExists
	}

	if err := u.deps.Password.Validate(params.Password); err != nil {
		return nil, err
	}

	secret, err := u.deps.Password.Hash(params.Password)
	if err != nil {
		return nil, err
	}

	user := models.UserDTO{
		Email:  params.Email,
		Secret: string(secret),
		Name:   params.Name,
	}

	result, err := u.deps.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	u.deps.UserCache.Set(result.Id, *result, cache.Expiration{Ttl: userCacheTtl})

	return result, nil
}

func (u *userService) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	user, err := u.deps.UserRepo.GetUserByEmail(ctx, email)
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

	u.deps.UserCache.Set(user.Id, *user, cache.Expiration{Ttl: userCacheTtl})

	return jwt, nil
}

func (u *userService) ForgotPassword(ctx context.Context, userId string) error {
	user, err := u.getUserById(ctx, userId)
	if err != nil {
		return err
	}

	actionToken, err := u.deps.ActionTokenRepo.CreateActionToken(
		ctx,
		models.ActionTokenDTO{
			UserId:     user.Id,
			Value:      uuid.New().String(),
			Target:     models.ActionTokenTargetForgotPassword,
			Expiration: time.Now().Add(1 * time.Hour),
		},
	)
	if err != nil {
		return err
	}

	return u.deps.EventRepo.SendEmailForgotPassword(ctx, user.Email, actionToken.Value)
}

func (u *userService) ChangePasswordWithToken(ctx context.Context, userId, actionToken, newPassword string) error {
	user, err := u.getUserById(ctx, userId)
	if err != nil {
		return err
	}

	code, err := u.deps.ActionTokenRepo.PopActionToken(ctx, userId, actionToken, models.ActionTokenTargetForgotPassword)
	if err != nil {
		return err
	}
	if code == nil {
		return fmt.Errorf("wrong user access code")
	}

	return u.updatePassword(ctx, *user, newPassword)
}

func (u *userService) ChangePassword(ctx context.Context, userId, oldPassword, newPassword string) error {
	user, err := u.getUserById(ctx, userId)
	if err != nil {
		return err
	}

	if !u.deps.Password.Compare(oldPassword, user.Secret) {
		return ErrUserWrongPassword
	}

	return u.updatePassword(ctx, *user, newPassword)
}

func (u *userService) updatePassword(ctx context.Context, user models.UserDTO, newPassword string) error {
	if err := u.deps.Password.Validate(newPassword); err != nil {
		return ErrUserBadPassword
	}

	u.deps.UserCache.Del(user.Id)

	newSecret, err := u.deps.Password.Hash(newPassword)
	if err != nil {
		return err
	}

	return u.deps.UserRepo.UpdateUser(ctx, user.Id, models.UserUpdateDTO{
		Secret: newSecret,
		Name:   user.Name,
	})
}

func (u *userService) getUserById(ctx context.Context, userId string) (*models.UserDTO, error) {
	if user, ok := u.deps.UserCache.GetEx(userId, cache.Expiration{Ttl: userCacheTtl}); ok {
		return &user, nil
	}

	user, err := u.deps.UserRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotExists
	}

	u.deps.UserCache.Set(user.Id, *user, cache.Expiration{Ttl: userCacheTtl})

	return user, nil
}

func (u *userService) ValidateToken(ctx context.Context, tokenStr string) (*models.UserDTO, error) {
	if userId, ok := u.deps.JwtCache.Get(tokenStr); ok {
		return u.getUserById(ctx, userId)
	}

	payload, err := u.deps.Jwt.Parse(tokenStr)
	if err != nil {
		return nil, ErrUserWrongToken
	}

	user, err := u.getUserById(ctx, payload.UserId)
	if err != nil {
		return nil, err
	}

	u.deps.JwtCache.Set(tokenStr, payload.UserId, cache.Expiration{ExpiresAt: payload.ExpiresAt.Time})

	return user, nil
}
