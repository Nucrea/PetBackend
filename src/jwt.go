package src

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JwtPayload struct {
	jwt.RegisteredClaims
	UserId string `json:"userId"`
}

type JwtUtil interface {
	Create(user UserDTO) (string, error)
	Parse(tokenStr string) (JwtPayload, error)
}

func NewJwtUtil(privateKey string) JwtUtil {
	return &jwtUtil{
		privateKey: privateKey,
	}
}

type jwtUtil struct {
	privateKey string
}

func (j *jwtUtil) Create(user UserDTO) (string, error) {
	payload := JwtPayload{UserId: user.Id}
	token := jwt.NewWithClaims(&jwt.SigningMethodHMAC{}, payload)
	tokenStr, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (j *jwtUtil) Parse(tokenStr string) (JwtPayload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, JwtPayload{}, func(t *jwt.Token) (interface{}, error) {
		return j.privateKey, nil
	})
	if err != nil {
		return JwtPayload{}, err
	}

	if payload, ok := token.Claims.(JwtPayload); ok {
		return payload, nil
	}

	return JwtPayload{}, fmt.Errorf("cant get payload")
}
