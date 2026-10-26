package order

// CreateOrderRequest is the body of POST /orders.
type CreateOrderRequest struct {
	ShippingName    string `json:"shipping_name"`
	ShippingPhone   string `json:"shipping_phone"`
	ShippingAddress string `json:"shipping_address"`
	Note            string `json:"note"`
	CouponCode      string `json:"coupon_code"`
	PaymentMethod   int16  `json:"payment_method"`
}

// OrderItemResponse is a snapshot line in an order.
type OrderItemResponse struct {
	ProductID    *string `json:"product_id,omitempty"`
	ProductTitle string  `json:"product_title"`
	ProductCover string  `json:"product_cover"`
	UnitPrice    float64 `json:"unit_price"`
	Quantity     int32   `json:"quantity"`
	Subtotal     float64 `json:"subtotal"`
}

// OrderResponse is the public representation of an order. List responses omit
// Items; detail responses include them.
type OrderResponse struct {
	ID              string              `json:"id"`
	Status          int16               `json:"status"`
	PaymentMethod   int16               `json:"payment_method"`
	PaymentStatus   int16               `json:"payment_status"`
	Subtotal        float64             `json:"subtotal"`
	DiscountAmount  float64             `json:"discount_amount"`
	ShippingFee     float64             `json:"shipping_fee"`
	Total           float64             `json:"total"`
	ShippingName    string              `json:"shipping_name"`
	ShippingPhone   string              `json:"shipping_phone"`
	ShippingAddress string              `json:"shipping_address"`
	Note            string              `json:"note"`
	CouponCode      string              `json:"coupon_code,omitempty"`
	PaidAt          *string             `json:"paid_at,omitempty"`
	CreateTime      string              `json:"create_time"`
	UpdateTime      string              `json:"update_time"`
	Items           []OrderItemResponse `json:"items,omitempty"`
}

// ListFilter controls order listing. Query values stay numeric.
type ListFilter struct {
	Status        *int16
	PaymentStatus *int16
	Page          int
	PerPage       int
}

// ListResult is a page of orders.
type ListResult struct {
	Items   []OrderResponse `json:"items"`
	Total   int64           `json:"total"`
	Page    int             `json:"page"`
	PerPage int             `json:"per_page"`
}

// UpdateStatusRequest is the body of PATCH /admin/orders/{id}/status.
type UpdateStatusRequest struct {
	Status int16 `json:"status"`
}

// UpdatePaymentRequest is the body of PATCH /admin/orders/{id}/payment.
type UpdatePaymentRequest struct {
	PaymentStatus int16 `json:"payment_status"`
}
