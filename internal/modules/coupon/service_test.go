package coupon

import (
	"math"
	"testing"
)

func floatPtr(v float64) *float64 { return &v }

func TestDiscount(t *testing.T) {
	tests := []struct {
		name     string
		coupon   *Coupon
		subtotal float64
		want     float64
	}{
		{"fixed", &Coupon{Type: CouponTypeFixed, Value: 10}, 50, 10},
		{"fixed clamps to subtotal", &Coupon{Type: CouponTypeFixed, Value: 10}, 5, 5},
		{"percent", &Coupon{Type: CouponTypePercent, Value: 20}, 250, 50},
		{"percent capped", &Coupon{Type: CouponTypePercent, Value: 20, MaxDiscount: floatPtr(30)}, 250, 30},
		{"percent under cap", &Coupon{Type: CouponTypePercent, Value: 10, MaxDiscount: floatPtr(30)}, 100, 10},
		{"percent rounds half up", &Coupon{Type: CouponTypePercent, Value: 15}, 99.99, 15},
		{"zero subtotal", &Coupon{Type: CouponTypeFixed, Value: 10}, 0, 0},
		{"nil coupon", nil, 100, 0},
		{"unknown type", &Coupon{Type: 9, Value: 10}, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Discount(tt.coupon, tt.subtotal)
			if math.Abs(got-tt.want) > 0.001 {
				t.Fatalf("Discount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeCode(t *testing.T) {
	if got := normalizeCode("  summer20 "); got != "SUMMER20" {
		t.Fatalf("normalizeCode() = %q", got)
	}
}
