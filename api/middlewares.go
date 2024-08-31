package api

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const (
	claimsId  = ctxKey("id")
	expiresAt = ctxKey("expires_at")
)

func (h *handler) authMiddleware(next http.Handler) http.Handler {
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
		claims, err := h.tokenMaker.VerifyToken(tokenStr)
		if err != nil {
			respondJSON(w, http.StatusUnauthorized, map[string]string{
				"message": "Error verifying token",
				"error":   err.Error(),
			})
			return
		}
		// Add claims to request context
		ctx := context.WithValue(r.Context(), claimsId, claims.RegisteredClaims.ID)
		ctx = context.WithValue(ctx, expiresAt, claims.RegisteredClaims.ExpiresAt.Time.Format("Mon Jan _2 15:04:05 2006"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
