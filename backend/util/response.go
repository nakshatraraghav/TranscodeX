package util

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes the given data structure as JSON to the response writer.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// WriteError writes an error message as JSON to the response writer.
func WriteError(w http.ResponseWriter, status int, err any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errorResponse := map[string]any{"error": err}
	return json.NewEncoder(w).Encode(errorResponse)
}
