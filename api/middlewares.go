package api

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	userId    = ctxKey("user_id")
	expiresAt = ctxKey("expires_at")
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondJSON(w, http.StatusUnauthorized, map[string]string{
				"message": "Authorization key is empty",
			})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			respondJSON(w, http.StatusUnauthorized, map[string]string{
				"message": "Invalid token format",
			})
			return
		}
		// Parse and validate the token
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil || !token.Valid {
			respondJSON(w, http.StatusUnauthorized, map[string]string{
				"message": err.Error(),
			})
			return
		}
		// Add claims to request context
		ctx := context.WithValue(r.Context(), userId, claims.Subject)
		ctx = context.WithValue(ctx, expiresAt, claims.ExpiresAt.Format("Mon Jan _2 15:04:05 2006"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
