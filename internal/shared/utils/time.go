package utils

import (
	"strings"
	"time"
)

// ParseTimePtr parses an optional RFC3339 timestamp. Blank and nil inputs
// return a nil timestamp.
func ParseTimePtr(raw *string) (*time.Time, error) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(*raw))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// FormatTimePtr formats an optional timestamp as RFC3339.
func FormatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
