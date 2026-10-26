package order

import "testing"

func TestOrderStatusTransitions(t *testing.T) {
	tests := []struct {
		from OrderStatus
		to   OrderStatus
		want bool
	}{
		{OrderStatusPending, OrderStatusConfirmed, true},
		{OrderStatusPending, OrderStatusCancelled, true},
		{OrderStatusConfirmed, OrderStatusProcessing, true},
		{OrderStatusConfirmed, OrderStatusCancelled, true},
		{OrderStatusProcessing, OrderStatusShipping, true},
		{OrderStatusProcessing, OrderStatusCancelled, true},
		{OrderStatusShipping, OrderStatusCompleted, true},
		{OrderStatusCompleted, OrderStatusRefunded, true},
		{OrderStatusPending, OrderStatusCompleted, false},
		{OrderStatusShipping, OrderStatusCancelled, false},
		{OrderStatusCancelled, OrderStatusConfirmed, false},
		{OrderStatusRefunded, OrderStatusPending, false},
	}

	for _, tt := range tests {
		if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
			t.Fatalf("%s -> %s = %v, want %v", tt.from, tt.to, got, tt.want)
		}
	}
}

func TestPaymentStatusTransitions(t *testing.T) {
	tests := []struct {
		from PaymentStatus
		to   PaymentStatus
		want bool
	}{
		{PaymentStatusUnpaid, PaymentStatusPaid, true},
		{PaymentStatusUnpaid, PaymentStatusFailed, true},
		{PaymentStatusFailed, PaymentStatusPaid, true},
		{PaymentStatusPaid, PaymentStatusRefunded, true},
		{PaymentStatusRefunded, PaymentStatusPaid, false},
		{PaymentStatusPaid, PaymentStatusUnpaid, false},
	}

	for _, tt := range tests {
		if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
			t.Fatalf("%s -> %s = %v, want %v", tt.from, tt.to, got, tt.want)
		}
	}
}

func TestEnumsValid(t *testing.T) {
	if !PaymentMethodCOD.Valid() {
		t.Fatal("cod should be valid")
	}
	if PaymentMethod(1).Valid() {
		t.Fatal("non-cod payment method should be invalid")
	}
	if PaymentMethod(9).Valid() {
		t.Fatal("unknown payment method should be invalid")
	}
	if !OrderStatusPending.Valid() || OrderStatus(9).Valid() {
		t.Fatal("order status validation is wrong")
	}
	if !PaymentStatusUnpaid.Valid() || PaymentStatus(9).Valid() {
		t.Fatal("payment status validation is wrong")
	}
}
