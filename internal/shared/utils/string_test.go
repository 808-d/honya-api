package utils

import "testing"

func TestOptionalString(t *testing.T) {
	if got := OptionalString(""); got != nil {
		t.Fatalf("OptionalString(\"\") = %v, want nil", got)
	}
	if got := OptionalString("hi"); got == nil || *got != "hi" {
		t.Fatalf("OptionalString(\"hi\") = %v", got)
	}
}
