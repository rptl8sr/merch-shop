package httputil

import (
	"encoding/json"
	"net/http"

	"merch-shop/internal/api"
)

func RespondError(w http.ResponseWriter, errMsg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(api.ErrorResponse{Errors: &errMsg})
}

func RespondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		RespondError(w, "Failed to write response", http.StatusInternalServerError)
	}
}
