package auth

import (
	"net/http"

	"github.com/passwordhash/jwt-test-task/internal/handler/api/v1/middleware"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/auth/tokens", h.tokens)
	mux.HandleFunc("/api/v1/auth/refresh", h.refresh)
	mux.Handle("/api/v1/auth/me", middleware.Identity(http.HandlerFunc(h.identify), h.tokensProvider))
	mux.Handle("/api/v1/auth/logout", middleware.Identity(http.HandlerFunc(h.logout), h.tokensProvider))
}
