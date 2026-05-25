package crypt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func getJwtKey() []byte {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		key = "fallback-dev-secret"
	}
	return []byte(key)
}

func CreateToken(userId string) (string, error) {
	claims := CustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "project-timegate",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJwtKey())
}

func VerifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return getJwtKey(), nil
		})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
