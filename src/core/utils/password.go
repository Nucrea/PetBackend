package utils

import (
	"backend/src/charsets"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordUtil interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
	Validate(password string) error
}

func NewPasswordUtil() PasswordUtil {
	specialChars := `!@#$%^&*()_-+={[}]|\:;"'<,>.?/`
	return &passwordUtil{
		charsetSpecialChars: charsets.NewCharsetFromString(specialChars),
	}
}

type passwordUtil struct {
	charsetSpecialChars charsets.Charset
}

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

	charsetUpper := charsets.GetCharset(charsets.CharsetTypeLettersUpper)
	charsetLower := charsets.GetCharset(charsets.CharsetTypeLettersLower)

	lowercaseLettersCount := 0
	uppercaseLettersCount := 0
	specialCharsCount := 0
	for _, v := range password {
		if b.charsetSpecialChars.TestRune(v) {
			specialCharsCount++
			continue
		}

		if charsetUpper.TestRune(v) {
			uppercaseLettersCount++
			continue
		}

		if charsetLower.TestRune(v) {
			lowercaseLettersCount++
			continue
		}
	}

	if lowercaseLettersCount == 0 {
		return fmt.Errorf("password must contain at least 1 lowercase letter")
	}
	if uppercaseLettersCount == 0 {
		return fmt.Errorf("password must contain at least 1 uppercase letter")
	}
	if specialCharsCount == 0 {
		return fmt.Errorf("password must contain at least 1 special character")
	}

	return nil
}
