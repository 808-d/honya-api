package response

import (
	"encoding/json"
	"net/http"
)

// Envelope is the standard JSON shape for all API responses.
// Every field is emitted in snake_case (see json tags).
type Envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func write(w http.ResponseWriter, status int, payload Envelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// OK sends a 200 response with data.
func OK(w http.ResponseWriter, message string, data any) {
	write(w, http.StatusOK, Envelope{Success: true, Message: message, Data: data})
}

// Created sends a 201 response with data.
func Created(w http.ResponseWriter, message string, data any) {
	write(w, http.StatusCreated, Envelope{Success: true, Message: message, Data: data})
}

// Error sends an error response with the given HTTP status and error code.
func Error(w http.ResponseWriter, status int, code, message string) {
	write(w, status, Envelope{Success: false, Message: message, Error: code})
}

// BadRequest is a 400 helper.
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, "bad_request", message)
}

// Unauthorized is a 401 helper.
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, "unauthorized", message)
}

// Internal is a 500 helper.
func Internal(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "internal_error", "internal server error")
}
