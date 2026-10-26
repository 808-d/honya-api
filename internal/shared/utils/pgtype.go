package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// TextPtr converts a nullable text column into a pointer.
func TextPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

// TextFromPtr converts a pointer into a nullable text parameter.
func TextFromPtr(s *string) pgtype.Text {
	if s == nil || *s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// TextFromString converts a string into a nullable text parameter, treating the
// empty string as NULL.
func TextFromString(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

// TimestamptzPtr converts a nullable timestamptz column into a pointer.
func TimestamptzPtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
}

// TimestamptzFromPtr converts a pointer into a nullable timestamptz parameter.
func TimestamptzFromPtr(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// Int2FromPtr converts a pointer into a nullable int2 parameter.
func Int2FromPtr(v *int16) pgtype.Int2 {
	if v == nil {
		return pgtype.Int2{}
	}
	return pgtype.Int2{Int16: *v, Valid: true}
}

// Int8Ptr converts a nullable int8 column into a pointer.
func Int8Ptr(v pgtype.Int8) *int64 {
	if !v.Valid {
		return nil
	}
	n := v.Int64
	return &n
}

// Int8FromPtr converts a pointer into a nullable int8 parameter.
func Int8FromPtr(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}

// Int4Ptr converts a nullable int4 column into a pointer.
func Int4Ptr(v pgtype.Int4) *int32 {
	if !v.Valid {
		return nil
	}
	n := v.Int32
	return &n
}

// Int4FromPtr converts a pointer into a nullable int4 parameter.
func Int4FromPtr(v *int32) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *v, Valid: true}
}

// Float8FromPtr converts a pointer into a nullable float8 parameter.
func Float8FromPtr(v *float64) pgtype.Float8 {
	if v == nil {
		return pgtype.Float8{}
	}
	return pgtype.Float8{Float64: *v, Valid: true}
}

// NumericPtr converts a nullable numeric column into a float pointer.
func NumericPtr(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	v, err := n.Float64Value()
	if err != nil || !v.Valid {
		return nil
	}
	f := v.Float64
	return &f
}

// BoolFromPtr converts a pointer into a nullable bool parameter.
func BoolFromPtr(v *bool) pgtype.Bool {
	if v == nil {
		return pgtype.Bool{}
	}
	return pgtype.Bool{Bool: *v, Valid: true}
}
