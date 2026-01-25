package utils

import (
	"context"
	"errors"

	"os"
	"time"
	
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id, role string) (string, error) {


	myToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("SECRET")
	if secret == "" {
		return "", errors.New("thers no secret key found")
	}
	tokenString, err := myToken.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("couldnt generate a token ")
	}
	return tokenString, nil

}

type contextKey string

const claimsKey contextKey = "jwtClaims"

func AddClaimsToContext(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}
