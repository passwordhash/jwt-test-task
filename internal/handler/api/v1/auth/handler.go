package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/passwordhash/jwt-test-task/internal/handler/api/v1/middleware"
	"github.com/passwordhash/jwt-test-task/internal/handler/api/v1/response"
)

type TokensProvider interface {
	GetPair(ctx context.Context, id, remoteAddr, userAgent string) (access, refresh string, err error)
	UserIDByToken(ctx context.Context, token string) (string, error)
}

type TokenRevoker interface {
	RevokeRefreshToken(ctx context.Context, userID, userAgent string) error
}

type Handler struct {
	tokensProvider TokensProvider
	tokenRevoker   TokenRevoker
}

func New(
	tokensProvider TokensProvider,
	tokenRevoker TokenRevoker,
) *Handler {
	return &Handler{
		tokensProvider: tokensProvider,
		tokenRevoker:   tokenRevoker,
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
		// TODO: add proper handling
		http.Error(w, "id and User-Agent are required", http.StatusBadRequest)
		return
	}

	access, refresh, err := h.tokensProvider.GetPair(r.Context(), id, r.RemoteAddr, userAgent)
	if err != nil {
		// TODO: add proper handling
		http.Error(w, fmt.Sprintf("failed to get tokens: %v", err), http.StatusInternalServerError)
		return
	}

	response.OK(w, tokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	if !response.ValidateMethod(r, w, http.MethodPost) {
		return
	}

	_, _ = fmt.Fprintf(w, "refresh endpoint")
}

func (h *Handler) identify(w http.ResponseWriter, r *http.Request) {
	if !response.ValidateMethod(r, w, http.MethodGet) {
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		response.Unauthorized(w, "Unable to identify user")
		return
	}

	response.OK(w, userIDResponse{
		UserID: userID,
	})
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if !response.ValidateMethod(r, w, http.MethodPost) {
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		response.Unauthorized(w, "Unable to identify user")
		return
	}

	err := h.tokenRevoker.RevokeRefreshToken(r.Context(), userID, r.Header.Get("User-Agent"))
	if err != nil {
		// TODO: add proper handling
		http.Error(w, fmt.Sprintf("failed to revoke token: %v", err), http.StatusInternalServerError)
		return
	}

	response.OK(w, "Token revoked successfully")
}
