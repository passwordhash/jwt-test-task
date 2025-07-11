package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/passwordhash/jwt-test-task/internal/handler/api/v1/response"
)

type TokensProvider interface {
	GetPair(ctx context.Context, id, remoteAddr, userAgent string) (access, refresh string, err error)
}

type Handler struct {
	tokensProvider TokensProvider
}

func New(
	tokensProvider TokensProvider,
) *Handler {
	return &Handler{
		tokensProvider: tokensProvider,
	}
}

func (h *Handler) token(w http.ResponseWriter, r *http.Request) {
	if !response.ValidateMethod(r, w, http.MethodPost) {
		return
	}

	id, ok := response.ValidateIDQueryParam(r, w)
	if !ok {
		return
	}
	userAgent := r.Header.Get("User-Agent")

	if id == "" || userAgent == "" {
		http.Error(w, "id and User-Agent are required", http.StatusBadRequest)
		return
	}

	access, refresh, err := h.tokensProvider.GetPair(r.Context(), id, r.RemoteAddr, userAgent)
	if err != nil {
		// TODO: add proper handling
		http.Error(w, fmt.Sprintf("failed to get tokens: %v", err), http.StatusInternalServerError)
		return
	}

	// TODO: temp
	resp := map[string]string{
		"access_token":  access,
		"refresh_token": refresh,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	respBody, _ := json.Marshal(resp)

	_, _ = fmt.Fprint(w, string(respBody))
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	if !response.ValidateMethod(r, w, http.MethodPost) {
		return
	}

	_, _ = fmt.Fprintf(w, "refresh endpoint")
}

func (h *Handler) idByToken(w http.ResponseWriter, r *http.Request) {
	if !response.ValidateMethod(r, w, http.MethodGet) {
		return
	}

	_, _ = fmt.Fprintf(w, "idByToken endpoint")
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if !response.ValidateMethod(r, w, http.MethodPost) {
		return
	}

	_, _ = fmt.Fprintf(w, "logout endpoint")
}
