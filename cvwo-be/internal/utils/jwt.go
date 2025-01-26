package utils

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("your-secret-key")

// CreateToken generates a new JWT token
func CreateToken(userID uint, name string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(userID), 10), // UserID as subject
		Issuer:    "MaiSpace",  // Issuer
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Expiry time (15 minutes)
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyToken parses and verifies the JWT token
func VerifyToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func CreateRefreshToken(userID uint) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   string(rune(userID)), // UserID as subject
		Issuer:    "MaiSpace",           // Issuer
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 1 week expiration
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Define a custom type for context keys
type contextKey string

const claimsKey contextKey = "claims"

// Middleware to check if the JWT token is valid
func IsAuthorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// Remove 'Bearer ' prefix if it exists
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		// Verify the token
		claims, err := VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Attach claims to context if you need to use it later
		r = r.WithContext(context.WithValue(r.Context(), claimsKey, claims))

		next(w, r)
	}
}