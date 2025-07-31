// Package middleware deals with middlewares
package middleware

import (
	"chatter/server/internal/user"
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "user"

func Auth(publicKey *rsa.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenStr string
			fmt.Printf("r.URL.Header().Get(\"Connection\"): %v\n", r.Header.Get("Connection"))
			fmt.Printf("r.URL.Header().Get(\"Upgrade\"): %v\n", r.Header.Get("Upgrade"))
			if isWebSocket(r) {
				log.Println("Websocket connection")
				tokenStr = r.URL.Query().Get("token")
			} else {
				tokenStr = extractTokenFromHeader(r)
			}
			if tokenStr == "" {
				http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
				return
			}

			claims, err := parseJWT(tokenStr, publicKey)
			if err != nil {
				log.Printf("middleware: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

func parseJWT(tokenStr string, publicKey *rsa.PublicKey) (*user.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &user.CustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
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

	fmt.Printf("claims: %v\n", claims)

	return claims, nil
}

func isWebSocket(r *http.Request) bool {
	connectionStr := strings.ToLower(r.Header.Get("Connection"))
	upgradeStr := strings.ToLower(r.Header.Get("Upgrade"))

	return strings.Contains(connectionStr, "upgrade") && strings.Contains(upgradeStr, "websocket")
}
