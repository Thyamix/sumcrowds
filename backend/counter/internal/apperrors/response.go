package apperrors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendError(w http.ResponseWriter, appError *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appError.StatusCode)

	jsonErr := json.NewEncoder(w).Encode(map[string]string{
		"error": appError.Public,
	})

	if jsonErr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	fmt.Print(appError.Internal)
}


