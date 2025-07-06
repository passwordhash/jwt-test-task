package auth

import (
	"fmt"
	"net/http"

	"github.com/passwordhash/jwt-test-task/internal/handler/response"
)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) token(w http.ResponseWriter, r *http.Request) {
	if !response.ValidateMethod(r, w, http.MethodPost) {
		return
	}

	fmt.Fprintf(w, "token endpoint")
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
