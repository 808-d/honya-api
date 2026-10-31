package order

import (
	"strings"
	"testing"
	"time"
)

func TestBuildOrderConfirmationEmail(t *testing.T) {
	contact := &UserContact{
		Username:    "reader",
		DisplayName: "",
		Email:       "reader@example.com",
	}
	o := &Order{
		PublicID:        "order-public-id",
		Subtotal:        100,
		DiscountAmount:  10,
		ShippingFee:     10,
		Total:           100,
		ShippingName:    "Reader",
		ShippingPhone:   "0123456789",
		ShippingAddress: "1 Book Street",
		Note:            "Leave at the door",
		CreateTime:      time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC),
		Items: []OrderItem{{
			ProductTitle: "Manga Vol. 1",
			ProductCover: "data:image/png;base64,VERY_LONG_COVER",
			UnitPrice:    50,
			Quantity:     2,
			Subtotal:     100,
		}},
	}

	msg, err := buildOrderConfirmationEmail(contact, o, "https://honya.example/")
	if err != nil {
		t.Fatalf("buildOrderConfirmationEmail() error = %v", err)
	}

	if msg.To != contact.Email {
		t.Fatalf("To = %q, want %q", msg.To, contact.Email)
	}
	if msg.Subject != "Your Honya order is confirmed" {
		t.Fatalf("Subject = %q", msg.Subject)
	}

	for _, body := range []string{msg.Text, msg.HTML} {
		for _, want := range []string{
			"reader",
			"order-public-id",
			"Manga Vol. 1",
			"¥100.00",
			"https://honya.example/orders/order-public-id",
			"Leave at the door",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("email body missing %q", want)
			}
		}
		if strings.Contains(body, "VERY_LONG_COVER") {
			t.Fatal("email body contains product cover")
		}
	}
}

func TestBuildOrderConfirmationEmailFallsBackToUsername(t *testing.T) {
	msg, err := buildOrderConfirmationEmail(
		&UserContact{Username: "reader", Email: "reader@example.com"},
		&Order{PublicID: "id", CreateTime: time.Now()},
		"",
	)
	if err != nil {
		t.Fatalf("buildOrderConfirmationEmail() error = %v", err)
	}
	if !strings.Contains(msg.Text, "Hi reader,") {
		t.Fatalf("text = %q", msg.Text)
	}
	if !strings.Contains(msg.Text, "/orders/id") {
		t.Fatalf("text = %q", msg.Text)
	}
}
