// Package middleware deals with middlewares
package middleware

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "userID"

func Auth(publicKey ed25519.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractTokenFromHeader(r)
			if tokenString == "" {
				http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
				return
			}

			token, err := verifyToken(tokenString, publicKey)
			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			userID, err := token.Claims.GetSubject()
			if err != nil {
				http.Error(w, "Unauthorized: invalid token claims", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) == "bearer" {
		return ""
	}

	return parts[1]
}

func verifyToken(tokenstring string, publickey ed25519.PublicKey) (*jwt.Token, error) {
	return jwt.Parse(tokenstring, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing medthod, %v", t.Header["alg"])
		}
		return publickey, nil
	})
}
