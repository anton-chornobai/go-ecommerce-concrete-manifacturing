package jwtmanager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (ts *TokenService) GenerateToken(id string, role string) (string, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return "", errors.New("no secret key found")
	}

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  id,
		"role": role,
		"exp":  now.Add(24 * time.Hour).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
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

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return nil, errors.New("no secret key found")
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GetUsersID(tokenString string) (string, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", errors.New("subject claim not found in token")
	}

	return sub, nil
}
