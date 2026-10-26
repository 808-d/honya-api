package order

import (
	"time"

	"github.com/808-d/honya-api/internal/modules/product"
)

// Order is a customer order and its immutable money snapshot.
type Order struct {
	ID              int64
	PublicID        string
	UserID          int64
	CouponCode      string
	Subtotal        float64
	DiscountAmount  float64
	ShippingFee     float64
	Total           float64
	Status          OrderStatus
	PaymentStatus   PaymentStatus
	PaymentMethod   PaymentMethod
	ShippingName    string
	ShippingPhone   string
	ShippingAddress string
	Note            string
	PaidAt          *time.Time
	CreateTime      time.Time
	UpdateTime      time.Time
	Items           []OrderItem
}

// OrderItem is a snapshot of a product at checkout time.
type OrderItem struct {
	ID                int64
	ProductInternalID *int64
	ProductID         *string
	ProductTitle      string
	ProductCover      string
	UnitPrice         float64
	Quantity          int32
	Subtotal          float64
}

// CheckoutLine is a locked cart line with its product snapshot.
type CheckoutLine struct {
	CartItemID      int64
	CartID          int64
	ProductID       int64
	ProductPublicID string
	ProductTitle    string
	ProductCover    string
	UnitPrice       float64
	Quantity        int32
	Stock           int32
	Status          product.ProductStatus
	IsDelete        bool
}

// CreateOrderWrite is a validated order ready to persist.
type CreateOrderWrite struct {
	UserID          int64
	CouponID        *int64
	Subtotal        float64
	DiscountAmount  float64
	ShippingFee     float64
	Total           float64
	Status          OrderStatus
	PaymentStatus   PaymentStatus
	PaymentMethod   PaymentMethod
	ShippingName    string
	ShippingPhone   string
	ShippingAddress string
	Note            string
}

// CreateOrderItemWrite is a snapshot line ready to persist.
type CreateOrderItemWrite struct {
	OrderID      int64
	ProductID    int64
	ProductTitle string
	ProductCover string
	UnitPrice    float64
	Quantity     int32
	Subtotal     float64
}
