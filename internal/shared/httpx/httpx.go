// Package httpx holds small request-parsing helpers shared by HTTP handlers.
package httpx

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/808-d/honya-api/internal/shared/apperror"
)

// Decode reads a JSON request body into dst.
func Decode(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

// QueryInt parses an optional integer query parameter, returning fallback
// when the value is empty.
func QueryInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		panic(apperror.BadRequest("invalid numeric query parameter"))
	}
	return v
}

// QueryInt16 parses an optional 16-bit integer query parameter.
func QueryInt16(raw string) *int16 {
	if raw == "" {
		return nil
	}
	v, err := strconv.ParseInt(raw, 10, 16)
	if err != nil {
		panic(apperror.BadRequest("invalid numeric query parameter"))
	}
	n := int16(v)
	return &n
}

// QueryFloat parses an optional float query parameter.
func QueryFloat(raw string) *float64 {
	if raw == "" {
		return nil
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		panic(apperror.BadRequest("invalid numeric query parameter"))
	}
	return &v
}
