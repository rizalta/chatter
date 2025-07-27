// Package middleware deals with middlewares
package middleware

import (
	"chatter/server/internal/user"
	"context"
	"crypto/ed25519"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userKey contextKey = "user"

func Auth(publicKey ed25519.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractTokenFromHeader(r)
			if tokenString == "" {
				http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
				return
			}

			claims, err := parseJWT(tokenString, publicKey)
			if err != nil {
				log.Printf("middleware: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userKey, claims)

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

func parseJWT(tokenStr string, publicKey ed25519.PublicKey) (*user.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &user.CustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("middleware: unexpected signing method, %v", t.Header["alg"])
		}

		return publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("unauthorized: invalid token, %v", err)
	}

	claims, ok := token.Claims.(*user.CustomClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized: invalid claims")
	}

	return claims, nil
}
