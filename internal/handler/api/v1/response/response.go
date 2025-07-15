package response

import (
	"encoding/json"
	"net/http"
)

func OK(w http.ResponseWriter, data interface{}) {
	jsonResponse(w, http.StatusOK, data)
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func ValidateMethod(r *http.Request, w http.ResponseWriter, expected string) bool {
	if r.Method != expected {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func ValidateIDQueryParam(r *http.Request, w http.ResponseWriter) (string, bool) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return "", false
	}

	return id, true
}
