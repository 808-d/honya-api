package coupon

// CouponInput is the body of the create and update endpoints.
type CouponInput struct {
	Code           string   `json:"code"`
	Type           int16    `json:"type"`
	Value          float64  `json:"value"`
	MinOrderAmount float64  `json:"min_order_amount"`
	MaxDiscount    *float64 `json:"max_discount"`
	StartAt        *string  `json:"start_at"`
	ExpireAt       *string  `json:"expire_at"`
	UsageLimit     *int32   `json:"usage_limit"`
	PerUserLimit   *int32   `json:"per_user_limit"`
	IsActive       *bool    `json:"is_active"`
}

// CouponResponse is the public representation of a coupon.
type CouponResponse struct {
	ID             string   `json:"id"`
	Code           string   `json:"code"`
	Type           int16    `json:"type"`
	Value          float64  `json:"value"`
	MinOrderAmount float64  `json:"min_order_amount"`
	MaxDiscount    *float64 `json:"max_discount"`
	StartAt        *string  `json:"start_at"`
	ExpireAt       *string  `json:"expire_at"`
	UsageLimit     *int32   `json:"usage_limit"`
	UsedCount      int32    `json:"used_count"`
	PerUserLimit   int32    `json:"per_user_limit"`
	IsActive       bool     `json:"is_active"`
	CreateTime     string   `json:"create_time"`
	UpdateTime     string   `json:"update_time"`
}

// ValidateRequest is the body of POST /coupons/validate.
type ValidateRequest struct {
	Code     string  `json:"code"`
	Subtotal float64 `json:"subtotal"`
}

// ValidateResponse is the discount preview for a coupon.
type ValidateResponse struct {
	Code           string  `json:"code"`
	Subtotal       float64 `json:"subtotal"`
	DiscountAmount float64 `json:"discount_amount"`
	Total          float64 `json:"total"`
}

// UsageResponse is one coupon redemption.
type UsageResponse struct {
	ID         int64  `json:"id"`
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	OrderID    string `json:"order_id"`
	CreateTime string `json:"create_time"`
}
