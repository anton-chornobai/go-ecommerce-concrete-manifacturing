package jwtmanager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenService struct {}

func NewTokenService() *TokenService {
	return  &TokenService{}
}

func (ts *TokenService) GenerateToken(id, role string) (string, error) {
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


func ValidateToken(stringToken string) (map[string]any, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return nil, errors.New("SECRET env variable missing")
	}

	token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot parse claims")
	}

	return map[string]interface{}(claims), nil
}
