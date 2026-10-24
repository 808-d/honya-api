package utils

import "testing"

func TestRound2(t *testing.T) {
	tests := []struct {
		in   float64
		want float64
	}{
		{10.126, 10.13},
		{10.124, 10.12},
		{14.9985, 15},
		{-1.235, -1.24},
		{0, 0},
	}
	for _, tt := range tests {
		if got := Round2(tt.in); got != tt.want {
			t.Fatalf("Round2(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
