package utils

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JwtPayload struct {
	UserId string `json:"userId"`
}

type JwtClaims struct {
	jwt.RegisteredClaims
	JwtPayload
}

type JwtUtil interface {
	Create(payload JwtPayload) (string, error)
	Parse(tokenStr string) (JwtClaims, error)
}

func NewJwtUtil(privateKey *rsa.PrivateKey) JwtUtil {
	return &jwtUtil{
		privateKey: privateKey,
	}
}

type jwtUtil struct {
	privateKey *rsa.PrivateKey
}

func (j *jwtUtil) Create(payload JwtPayload) (string, error) {
	claims := &JwtClaims{JwtPayload: payload}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (j *jwtUtil) Parse(tokenStr string) (JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return &j.privateKey.PublicKey, nil
	})
	if err != nil {
		return JwtClaims{}, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok {
		return *claims, nil
	}

	return JwtClaims{}, fmt.Errorf("cant get payload")
}
