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
