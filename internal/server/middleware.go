package server

import (
	"context"
	"net/http"
	"strings"

	"iss-model-backend/internal/utils"
)

type contextKey string

const userContextKey = contextKey("user")

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Missing authorization header", "")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid token format", "")
			return
		}

		claims, err := s.authService.ValidateToken(tokenString)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid token", err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
