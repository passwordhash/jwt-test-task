package auth

import "net/http"

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/auth/token", h.token)
	mux.HandleFunc("/api/v1/auth/refresh", h.refresh)
	mux.HandleFunc("/api/v1/auth/me", h.idByToken)  // TODO: protected endpoint
	mux.HandleFunc("/api/v1/auth/logout", h.logout) // TODO: protected endpoint
}
