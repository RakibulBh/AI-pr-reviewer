package json

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Envelope represents a standard API response structure
type Envelope struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// WriteJSON writes a JSON response to the HTTP response writer
func WriteJSON(w http.ResponseWriter, status int, envelope Envelope) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(envelope); err != nil {
		return fmt.Errorf("failed to encode JSON response: %w", err)
	}

	return nil
}

// WriteSuccessJSON writes a successful JSON response
func WriteSuccessJSON(w http.ResponseWriter, status int, message string, data any) error {
	envelope := Envelope{
		Error:   false,
		Message: message,
		Data:    data,
	}

	return WriteJSON(w, status, envelope)
}

// WriteErrorJSON writes an error JSON response
func WriteErrorJSON(w http.ResponseWriter, status int, message string) error {
	envelope := Envelope{
		Error:   true,
		Message: message,
		Data:    nil,
	}

	return WriteJSON(w, status, envelope)
}

// ReadJSON reads and parses JSON from HTTP request body
func ReadJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("failed to decode JSON request: %w", err)
	}

	return nil
}

// ReadJSONWithValidation reads JSON and validates content type
func ReadJSONWithValidation(r *http.Request, dst any) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("invalid content type: expected application/json, got %s", contentType)
	}

	return ReadJSON(r, dst)
}

// WriteJSONError is a convenience function for common error responses
func WriteJSONError(w http.ResponseWriter, err error) error {
	return WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
}

// WriteBadRequestJSON writes a 400 Bad Request JSON response
func WriteBadRequestJSON(w http.ResponseWriter, message string) error {
	return WriteErrorJSON(w, http.StatusBadRequest, message)
}

// WriteNotFoundJSON writes a 404 Not Found JSON response
func WriteNotFoundJSON(w http.ResponseWriter, message string) error {
	return WriteErrorJSON(w, http.StatusNotFound, message)
}

// WriteUnauthorizedJSON writes a 401 Unauthorized JSON response
func WriteUnauthorizedJSON(w http.ResponseWriter, message string) error {
	return WriteErrorJSON(w, http.StatusUnauthorized, message)
}
