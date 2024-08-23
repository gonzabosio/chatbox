package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey}
}

func (m *JWTMaker) CreateToken(id, username string, duration time.Duration) (string, *CustomClaims, error) {
	claims, err := NewCustomClaims(id, username, duration)
	if err != nil {
		return "", nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", nil, err
	}
	return tokenStr, claims, nil
}

func (m *JWTMaker) VerifyToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		//verify signing method
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return []byte(m.secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
