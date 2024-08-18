package middleware

import (
	"context"
	"github.com/olzzhas/narxozer/graph"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextUserID   contextKey = "user_id"
	ContextUserRole contextKey = "role"
)

func AuthMiddleware(manager *graph.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := manager.Verify(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Добавляем UserID и Role в контекст запроса
			ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextUserRole, claims.Role)

			// Передаем управление следующему обработчику
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) int64 {
	if userID, ok := ctx.Value(ContextUserID).(int64); ok {
		return userID
	}
	return 0
}

func GetUserRoleFromContext(ctx context.Context) string {
	if role, ok := ctx.Value(ContextUserRole).(string); ok {
		return role
	}
	return ""
}
