package order

import "testing"

func TestTotals(t *testing.T) {
	tests := []struct {
		name         string
		percent      float64
		subtotal     float64
		discount     float64
		wantShipping float64
		wantTotal    float64
	}{
		{"ten percent", 10, 100, 10, 10, 100},
		{"rounds shipping", 10, 99.99, 0, 10, 109.99},
		{"discount larger than subtotal is clamped by total", 10, 5, 5, 0.5, 0.5},
		{"zero percent", 0, 50, 0, 0, 50},
		{"never negative", 10, 10, 20, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{shippingFeePercent: tt.percent}
			shipping, total := svc.totals(tt.subtotal, tt.discount)
			if shipping != tt.wantShipping {
				t.Fatalf("shipping = %v, want %v", shipping, tt.wantShipping)
			}
			if total != tt.wantTotal {
				t.Fatalf("total = %v, want %v", total, tt.wantTotal)
			}
		})
	}
}

func TestNewServiceClampsShippingPercent(t *testing.T) {
	svc := NewService(nil, nil, -5, nil, "")
	if svc.shippingFeePercent != 0 {
		t.Fatalf("shippingFeePercent = %v, want 0", svc.shippingFeePercent)
	}
}

func TestCustomerCanCancel(t *testing.T) {
	tests := []struct {
		name    string
		status  OrderStatus
		payment PaymentStatus
		want    bool
	}{
		{"pending unpaid", OrderStatusPending, PaymentStatusUnpaid, true},
		{"confirmed unpaid", OrderStatusConfirmed, PaymentStatusUnpaid, true},
		{"processing failed payment", OrderStatusProcessing, PaymentStatusFailed, true},
		{"shipping unpaid", OrderStatusShipping, PaymentStatusUnpaid, false},
		{"pending paid", OrderStatusPending, PaymentStatusPaid, false},
		{"processing refunded", OrderStatusProcessing, PaymentStatusRefunded, false},
		{"completed unpaid", OrderStatusCompleted, PaymentStatusUnpaid, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := customerCanCancel(tt.status, tt.payment); got != tt.want {
				t.Fatalf("customerCanCancel(%v, %v) = %v, want %v", tt.status, tt.payment, got, tt.want)
			}
		})
	}
}

func TestCompletionMarksPaid(t *testing.T) {
	tests := []struct {
		payment PaymentStatus
		want    bool
	}{
		{PaymentStatusUnpaid, true},
		{PaymentStatusFailed, true},
		{PaymentStatusPaid, false},
		{PaymentStatusRefunded, false},
	}

	for _, tt := range tests {
		t.Run(tt.payment.String(), func(t *testing.T) {
			if got := completionMarksPaid(tt.payment); got != tt.want {
				t.Fatalf("completionMarksPaid(%v) = %v, want %v", tt.payment, got, tt.want)
			}
		})
	}
}
