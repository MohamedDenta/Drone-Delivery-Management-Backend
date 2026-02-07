package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SecretKey is the secret used to sign tokens. In production, load this from environment variables.
var SecretKey = []byte("super-secret-key-change-me")

// Claims represents the JWT claims
type Claims struct {
	Name     string `json:"name"`
	UserType string `json:"user_type"` // admin | enduser | drone
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(name, userType string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Name:     name,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   name, // Using name as subject for now
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// ValidateToken parses and validates the token, returning the claims
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
