package jwt

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrorNilToken   = errors.New("gets nil token")
	ErrInvalidToken = errors.New("invalid token")
)

func ValidToken(tokenString string) error {
	const op = "jwt.ValidToken"

	if tokenString == "" {
		return fmt.Errorf("%s: %w", op, ErrorNilToken)
	}

	secret := os.Getenv("APP_SECRET")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidToken)
		}
		return []byte(secret), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return fmt.Errorf("%s: %w", op, ErrTokenExpired)
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return fmt.Errorf("%s: %w", op, jwt.ErrTokenSignatureInvalid)
		case errors.Is(err, jwt.ErrTokenMalformed):
			return fmt.Errorf("%s: %w", op, jwt.ErrTokenMalformed)
		default:
			return fmt.Errorf("%s: %w", op, ErrInvalidToken)
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if _, exists := claims["app_id"]; !exists {
			return fmt.Errorf("%s: %w", op, ErrInvalidToken)
		}
		return nil
	}

	return fmt.Errorf("%s: %w", op, ErrInvalidToken)
}
