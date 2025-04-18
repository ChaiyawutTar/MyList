package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"github.com/golang-jwt/jwt/v4"
)

type JWTAuth struct {
	secretKey []byte
	expiry    time.Duration
}

func NewJWTAuth(secretKey string, expiry time.Duration) *JWTAuth {
	return &JWTAuth{
		secretKey: []byte(secretKey),
		expiry:    expiry,
	}
}

func (a *JWTAuth) GenerateToken(userID int) (string, error) {
	// Create claims
	claims := jwt.MapClaims{
		"user_id": strconv.Itoa(userID),
		"exp":     time.Now().Add(a.expiry).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	return token.SignedString(a.secretKey)
}

func (a *JWTAuth) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}