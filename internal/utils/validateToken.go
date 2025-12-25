package utils 

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func ValidateToken(stringToken string) (map[string]interface{}, error) {
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
