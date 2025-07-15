package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/passwordhash/jwt-test-task/internal/handler/api/v1/response"
)

type CtxKey string

const (
	UserIDKey CtxKey = "userID"
)

type IdentityProvider interface {
	UserIDByToken(ctx context.Context, token string) (string, error)
}

func Identity(next http.Handler, provider IdentityProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := tokenFromHeader(r)
		if err != nil {
			response.Unauthorized(w, fmt.Sprintf("failed to get token from header: %v", err))
			return
		}

		userID, err := provider.UserIDByToken(r.Context(), jwtToken)
		if err != nil || userID == "" {
			response.Unauthorized(w, fmt.Sprintf("failed to get user ID by token: %v", err))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func tokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) && len(authHeader) > len(prefix) {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return strings.TrimPrefix(authHeader, prefix), nil
}
