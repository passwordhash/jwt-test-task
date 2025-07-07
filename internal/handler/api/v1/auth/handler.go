package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/passwordhash/jwt-test-task/internal/handler/response"
)

type TokensProvider interface {
	GetPair(ctx context.Context, id, ip, userAgent string) (access, refresh string, err error)
}

type handler struct {
	tokensProvider TokensProvider
}

func New(
	tokensProvider TokensProvider,
) *handler {
	return &handler{
		tokensProvider: tokensProvider,
	}
}

func (h *handler) token(w http.ResponseWriter, r *http.Request) {
	if !response.ValidateMethod(r, w, http.MethodPost) {
		return
	}

	id, ok := response.ValidateIdQueryParam(r, w)
	if !ok {
		return
	}
	ip := getIP(r)
	userAgent := r.Header.Get("User-Agent")

	if id == "" || userAgent == "" {
		http.Error(w, "id and User-Agent are required", http.StatusBadRequest)
		return
	}

	access, refresh, err := h.tokensProvider.GetPair(r.Context(), id, ip, userAgent)
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

	fmt.Fprintf(w, string(respBody))
}

func (h *handler) refresh(w http.ResponseWriter, r *http.Request) {
	if !response.ValidateMethod(r, w, http.MethodPost) {
		return
	}

	fmt.Fprintf(w, "refresh endpoint")
}

func (h *handler) idByToken(w http.ResponseWriter, r *http.Request) {
	if !response.ValidateMethod(r, w, http.MethodGet) {
		return
	}

	fmt.Fprintf(w, "idByToken endpoint")
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	if !response.ValidateMethod(r, w, http.MethodPost) {
		return
	}

	fmt.Fprintf(w, "logout endpoint")
}

func getIP(r *http.Request) string {
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip = ip[:strings.LastIndex(ip, ":")]
	}
	return ip
}
