package response

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func OK(w http.ResponseWriter, data any) {
	jsonResponse(w, http.StatusOK, response{
		Success: true,
		Data:    data,
		Message: "",
	})
}

func Unauthorized(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusUnauthorized, response{
		Success: false,
		Data:    nil,
		Message: message,
	})
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

// ValidateMethod checks if the request method matches the expected method.
// If not, it responds with a 405 Method Not Allowed error.
//
// Note: in go 1.22, in url/http there is a new opportunity to use mux.HandleFunc("POST /api/v1/auth/tokens", h.token),
// but in this case, response type will be plain text, not json.
func ValidateMethod(r *http.Request, w http.ResponseWriter, expected string) bool {
	if r.Method != expected {
		jsonResponse(w, http.StatusMethodNotAllowed, response{
			Success: false,
			Data:    nil,
			Message: "Method not allowed",
		})
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
