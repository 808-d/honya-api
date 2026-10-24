package utils

import (
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestTextPtr(t *testing.T) {
	if got := TextPtr(pgtype.Text{}); got != nil {
		t.Fatalf("TextPtr(invalid) = %v, want nil", got)
	}
	if got := TextPtr(pgtype.Text{String: "hi", Valid: true}); got == nil || *got != "hi" {
		t.Fatalf("TextPtr(valid) = %v, want \"hi\"", got)
	}
}

func TestTextFromPtr(t *testing.T) {
	empty := ""
	value := "hi"
	tests := []struct {
		name string
		in   *string
		want pgtype.Text
	}{
		{"nil", nil, pgtype.Text{}},
		{"empty", &empty, pgtype.Text{}},
		{"value", &value, pgtype.Text{String: "hi", Valid: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TextFromPtr(tt.in); got != tt.want {
				t.Fatalf("TextFromPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextFromString(t *testing.T) {
	if got := TextFromString(""); got.Valid {
		t.Fatalf("TextFromString(\"\") = %v, want invalid", got)
	}
	if got := TextFromString("hi"); !got.Valid || got.String != "hi" {
		t.Fatalf("TextFromString(\"hi\") = %v", got)
	}
}

func TestTimestamptz(t *testing.T) {
	if got := TimestamptzPtr(pgtype.Timestamptz{}); got != nil {
		t.Fatalf("TimestamptzPtr(invalid) = %v, want nil", got)
	}
	if got := TimestamptzFromPtr(nil); got.Valid {
		t.Fatalf("TimestamptzFromPtr(nil) = %v, want invalid", got)
	}

	now := time.Now().UTC().Truncate(time.Second)
	got := TimestamptzPtr(TimestamptzFromPtr(&now))
	if got == nil || !got.Equal(now) {
		t.Fatalf("round trip = %v, want %v", got, now)
	}
}

func TestIntFromPtr(t *testing.T) {
	i16 := int16(3)
	if got := Int2FromPtr(&i16); !got.Valid || got.Int16 != i16 {
		t.Fatalf("Int2FromPtr() = %v", got)
	}
	if got := Int2FromPtr(nil); got.Valid {
		t.Fatalf("Int2FromPtr(nil) = %v", got)
	}

	i32 := int32(4)
	if got := Int4FromPtr(&i32); !got.Valid || got.Int32 != i32 {
		t.Fatalf("Int4FromPtr() = %v", got)
	}
	if got := Int4Ptr(Int4FromPtr(&i32)); got == nil || *got != i32 {
		t.Fatalf("Int4Ptr() = %v", got)
	}
	if got := Int4Ptr(pgtype.Int4{}); got != nil {
		t.Fatalf("Int4Ptr(invalid) = %v, want nil", got)
	}
}

func TestFloat8FromPtr(t *testing.T) {
	f := 12.5
	if got := Float8FromPtr(&f); !got.Valid || got.Float64 != f {
		t.Fatalf("Float8FromPtr() = %v", got)
	}
	if got := Float8FromPtr(nil); got.Valid {
		t.Fatalf("Float8FromPtr(nil) = %v", got)
	}
}

func TestNumericPtr(t *testing.T) {
	if got := NumericPtr(pgtype.Numeric{}); got != nil {
		t.Fatalf("NumericPtr(invalid) = %v, want nil", got)
	}

	n := pgtype.Numeric{Int: big.NewInt(1234), Exp: -2, Valid: true}
	if got := NumericPtr(n); got == nil || *got != 12.34 {
		t.Fatalf("NumericPtr() = %v, want 12.34", got)
	}
}

func TestBoolFromPtr(t *testing.T) {
	v := true
	if got := BoolFromPtr(&v); !got.Valid || !got.Bool {
		t.Fatalf("BoolFromPtr() = %v", got)
	}
	if got := BoolFromPtr(nil); got.Valid {
		t.Fatalf("BoolFromPtr(nil) = %v", got)
	}
}
