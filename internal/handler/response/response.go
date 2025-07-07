package response

import (
	"net/http"
)

func ValidateMethod(r *http.Request, w http.ResponseWriter, expected string) bool {
	if r.Method != expected {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func ValidateIdQueryParam(r *http.Request, w http.ResponseWriter) (string, bool) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return "", false
	}

	return id, true
}
