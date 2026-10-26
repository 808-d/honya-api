package coupon

import "time"

// Coupon is a discount code that can be applied to an order.
type Coupon struct {
	ID             int64
	PublicID       string
	Code           string
	Type           CouponType
	Value          float64
	MinOrderAmount float64
	MaxDiscount    *float64
	StartAt        *time.Time
	ExpireAt       *time.Time
	UsageLimit     *int32
	UsedCount      int32
	PerUserLimit   int32
	IsActive       bool
	CreateTime     time.Time
	UpdateTime     time.Time
}

// CouponWrite is a validated coupon payload ready to persist.
type CouponWrite struct {
	Code           string
	Type           CouponType
	Value          float64
	MinOrderAmount float64
	MaxDiscount    *float64
	StartAt        *time.Time
	ExpireAt       *time.Time
	UsageLimit     *int32
	PerUserLimit   int32
	IsActive       bool
}

// CouponUsage is one redemption of a coupon on an order.
type CouponUsage struct {
	ID         int64
	UserID     string
	Username   string
	OrderID    string
	CreateTime time.Time
}

// ListFilter controls coupon listing.
type ListFilter struct {
	Search    string
	Type      *int16
	IsActive  *bool
	ActiveNow bool
	Page      int
	PerPage   int
}

// ListResult is a page of coupons.
type ListResult struct {
	Items   []CouponResponse `json:"items"`
	Total   int64            `json:"total"`
	Page    int              `json:"page"`
	PerPage int              `json:"per_page"`
}

// UsageListResult is a page of coupon redemptions.
type UsageListResult struct {
	Items   []UsageResponse `json:"items"`
	Total   int64           `json:"total"`
	Page    int             `json:"page"`
	PerPage int             `json:"per_page"`
}
