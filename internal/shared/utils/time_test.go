package utils

import (
	"testing"
	"time"
)

func TestParseTimePtr(t *testing.T) {
	if got, err := ParseTimePtr(nil); err != nil || got != nil {
		t.Fatalf("ParseTimePtr(nil) = %v, %v", got, err)
	}

	blank := "   "
	if got, err := ParseTimePtr(&blank); err != nil || got != nil {
		t.Fatalf("ParseTimePtr(blank) = %v, %v", got, err)
	}

	raw := " 2026-01-02T03:04:05Z "
	got, err := ParseTimePtr(&raw)
	if err != nil {
		t.Fatalf("ParseTimePtr() error = %v", err)
	}
	want := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("ParseTimePtr() = %v, want %v", got, want)
	}

	bad := "not-a-time"
	if _, err := ParseTimePtr(&bad); err == nil {
		t.Fatal("ParseTimePtr(bad) error = nil, want error")
	}
}

func TestFormatTimePtr(t *testing.T) {
	if got := FormatTimePtr(nil); got != nil {
		t.Fatalf("FormatTimePtr(nil) = %v, want nil", got)
	}

	at := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	got := FormatTimePtr(&at)
	if got == nil || *got != "2026-01-02T03:04:05Z" {
		t.Fatalf("FormatTimePtr() = %v", got)
	}
}
