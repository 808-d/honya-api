package utils

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsUniqueViolation(t *testing.T) {
	if !IsUniqueViolation(&pgconn.PgError{Code: "23505"}) {
		t.Fatal("IsUniqueViolation(23505) = false, want true")
	}
	if IsUniqueViolation(&pgconn.PgError{Code: "23503"}) {
		t.Fatal("IsUniqueViolation(23503) = true, want false")
	}
	if IsUniqueViolation(errors.New("boom")) {
		t.Fatal("IsUniqueViolation(other) = true, want false")
	}
	if IsUniqueViolation(nil) {
		t.Fatal("IsUniqueViolation(nil) = true, want false")
	}
}
