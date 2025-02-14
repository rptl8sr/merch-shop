package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"merch-shop/internal/httputil"
	"merch-shop/pkg/jwt"
)

type contextKey string

const (
	ContextKeyUserID = contextKey("user_id")
	BearerAuthScopes = contextKey("scopes")
)

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				httputil.RespondError(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwt.ParseToken(secret, tokenStr)
			if err != nil {
				httputil.RespondError(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), BearerAuthScopes, []string{})
			ctx = context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value(ContextKeyUserID).(uint)
	if !ok {
		return 0, fmt.Errorf("user id not found in context")
	}

	return userID, nil
}
