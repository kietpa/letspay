package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTKey []byte

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func SetSecret(secret string) {
	JWTKey = []byte(secret)
}

func GenerateToken(userId int) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "letspay",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTKey)
}

func ValidateToken(tokenStr string) (Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		return JWTKey, nil
	})
	if err != nil || !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	return *claims, nil
}
