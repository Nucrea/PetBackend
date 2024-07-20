package src

import "golang.org/x/crypto/bcrypt"

type BCryptUtil interface {
	HashPassword(password string) (string, error)
	IsPasswordsEqual(password, hash string) bool
}

func NewBcrypt() BCryptUtil {
	return &bcryptImpl{}
}

type bcryptImpl struct{}

func (b *bcryptImpl) HashPassword(password string) (string, error) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), nil
}

func (b *bcryptImpl) IsPasswordsEqual(password, hash string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
