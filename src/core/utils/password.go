package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordUtil interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
	Validate(password string) error
}

func NewPasswordUtil() PasswordUtil {
	return &passwordUtil{}
}

type passwordUtil struct{}

func (b *passwordUtil) Hash(password string) (string, error) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), nil
}

func (b *passwordUtil) Compare(password, hash string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (b *passwordUtil) Validate(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must contain 8 or more characters")
	}
	return nil
}
