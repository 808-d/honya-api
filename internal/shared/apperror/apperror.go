package apperror

import "net/http"

// Error is an application error that carries the HTTP status and error code
// to return to the client. Handlers panic with it; the recover middleware
// translates it into a response.
type Error struct {
	Status  int
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func New(status int, code, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

func BadRequest(message string) *Error {
	return New(http.StatusBadRequest, "bad_request", message)
}

func Unauthorized(message string) *Error {
	return New(http.StatusUnauthorized, "unauthorized", message)
}

func Forbidden(message string) *Error {
	return New(http.StatusForbidden, "forbidden", message)
}

func NotFound(message string) *Error {
	return New(http.StatusNotFound, "not_found", message)
}

func Conflict(message string) *Error {
	return New(http.StatusConflict, "conflict", message)
}

func Internal(message string) *Error {
	return New(http.StatusInternalServerError, "internal_error", message)
}
