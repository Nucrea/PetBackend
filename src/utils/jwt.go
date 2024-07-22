package utils

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JwtPayload struct {
	UserId string `json:"userId"`
}

type jwtClaims struct {
	jwt.RegisteredClaims
	JwtPayload
}

type JwtUtil interface {
	Create(payload JwtPayload) (string, error)
	Parse(tokenStr string) (JwtPayload, error)
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
	claims := &jwtClaims{JwtPayload: payload}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (j *jwtUtil) Parse(tokenStr string) (JwtPayload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return &j.privateKey.PublicKey, nil
	})
	if err != nil {
		return JwtPayload{}, err
	}

	if claims, ok := token.Claims.(*jwtClaims); ok {
		return claims.JwtPayload, nil
	}

	return JwtPayload{}, fmt.Errorf("cant get payload")
}
