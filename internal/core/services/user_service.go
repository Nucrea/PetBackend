package services

import (
	"backend/internal/core/models"
	"backend/internal/core/repos"
	"backend/internal/core/utils"
	"backend/pkg/cache"
	"backend/pkg/logger"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotExists       = fmt.Errorf("no such user")
	ErrUserExists          = fmt.Errorf("user with this login already exists")
	ErrUserWrongPassword   = fmt.Errorf("wrong password")
	ErrUserWrongToken      = fmt.Errorf("bad user token")
	ErrUserBadPassword     = fmt.Errorf("password must contain at least 8 characters")
	ErrUserEmailUnverified = fmt.Errorf("user has not verified email yet")
	// ErrUserInternal = fmt.Errorf("unexpected error. contact tech support")
)

const (
	userCacheTtl = time.Hour
)

type UserService interface {
	CreateUser(ctx context.Context, params UserCreateParams) (*models.UserDTO, error)
	AuthenticateUser(ctx context.Context, login, password string) (string, error)
	ValidateAuthToken(ctx context.Context, tokenStr string) (*models.UserDTO, error)
	VerifyEmail(ctx context.Context, actionToken string) error

	SendEmailForgotPassword(ctx context.Context, userId string) error
	SendEmailVerifyUser(ctx context.Context, email string) error

	ChangePassword(ctx context.Context, userId, oldPassword, newPassword string) error
	ChangePasswordWithToken(ctx context.Context, actionToken, newPassword string) error
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
	Logger          logger.Logger
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
		Email:    strings.ToLower(params.Email),
		Secret:   string(secret),
		FullName: params.Name,
	}

	result, err := u.deps.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	if err := u.sendEmailVerifyUser(ctx, result.Id, user.Email); err != nil {
		u.deps.Logger.Error().Err(err).Msg("error occured on sending email")
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

	if !user.EmailVerified {
		return "", ErrUserEmailUnverified
	}

	payload := utils.JwtPayload{UserId: user.Id}
	jwt, err := u.deps.Jwt.Create(payload)
	if err != nil {
		return "", err
	}

	u.deps.UserCache.Set(user.Id, *user, cache.Expiration{Ttl: userCacheTtl})

	return jwt, nil
}

func (u *userService) VerifyEmail(ctx context.Context, actionToken string) error {
	token, err := u.deps.ActionTokenRepo.GetActionToken(ctx, actionToken, models.ActionTokenTargetVerifyEmail)
	if err != nil {
		return err
	}
	if token == nil {
		return fmt.Errorf("wrong action token")
	}

	if err := u.deps.UserRepo.SetUserEmailVerified(ctx, token.UserId); err != nil {
		return err
	}

	//TODO: log warnings somehow
	u.deps.ActionTokenRepo.DeleteActionToken(ctx, token.Id)
	return nil
}

func (u *userService) SendEmailForgotPassword(ctx context.Context, email string) error {
	// user, err := u.getUserById(ctx, userId)
	user, err := u.deps.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	actionToken, err := u.deps.ActionTokenRepo.CreateActionToken(
		ctx,
		models.ActionTokenDTO{
			UserId:     user.Id,
			Value:      uuid.New().String(),
			Target:     models.ActionTokenTargetRestorePassword,
			Expiration: time.Now().Add(15 * time.Minute),
		},
	)
	if err != nil {
		return err
	}

	return u.deps.EventRepo.SendEmailForgotPassword(ctx, user.Email, actionToken.Value)
}

func (u *userService) sendEmailVerifyUser(ctx context.Context, userId, email string) error {
	actionToken, err := u.deps.ActionTokenRepo.CreateActionToken(
		ctx,
		models.ActionTokenDTO{
			UserId:     userId,
			Value:      uuid.New().String(),
			Target:     models.ActionTokenTargetVerifyEmail,
			Expiration: time.Now().Add(1 * time.Hour),
		},
	)
	if err != nil {
		return err
	}

	return u.deps.EventRepo.SendEmailVerifyUser(ctx, email, actionToken.Value)
}

func (u *userService) SendEmailVerifyUser(ctx context.Context, email string) error {
	//user, err := u.getUserById(ctx, userId)
	user, err := u.deps.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("no such user")
	}
	if user.EmailVerified {
		return fmt.Errorf("user already verified")
	}

	return u.sendEmailVerifyUser(ctx, user.Id, user.Email)
}

func (u *userService) ChangePasswordWithToken(ctx context.Context, actionToken, newPassword string) error {
	token, err := u.deps.ActionTokenRepo.GetActionToken(ctx, actionToken, models.ActionTokenTargetRestorePassword)
	if err != nil {
		return err
	}
	if token == nil {
		return fmt.Errorf("wrong action token")
	}

	user, err := u.getUserById(ctx, token.UserId)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("no such user")
	}

	if err := u.updatePassword(ctx, *user, newPassword); err != nil {
		return err
	}

	//TODO: log warnings somehow
	u.deps.ActionTokenRepo.DeleteActionToken(ctx, token.Id)
	return nil
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

	if err = u.deps.UserRepo.UpdateUser(ctx, user.Id, models.UserUpdateDTO{
		Secret:   newSecret,
		FullName: user.FullName,
	}); err != nil {
		return err
	}

	if err := u.deps.EventRepo.SendEmailPasswordChanged(ctx, user.Email); err != nil {
		u.deps.Logger.Error().Err(err).Msg("error occured on sending email")
	}

	return nil
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

func (u *userService) ValidateAuthToken(ctx context.Context, tokenStr string) (*models.UserDTO, error) {
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
