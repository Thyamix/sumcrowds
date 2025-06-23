package apperrors

import (
	"encoding/json"
	"net/http"
)

func sendError(w http.ResponseWriter, appError *AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appError.StatusCode)

	jsonErr := json.NewEncoder(w).Encode(map[string]string{
		"error": appError.Public,
	})

	if jsonErr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
