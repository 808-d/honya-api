package dashboard

// ProductStats is the catalog summary.
type ProductStats struct {
	Total    int64 `json:"total"`
	Active   int64 `json:"active"`
	Draft    int64 `json:"draft"`
	Archived int64 `json:"archived"`
	LowStock int64 `json:"low_stock"`
}

// OrderStats is the order summary, including a breakdown by fulfillment state.
type OrderStats struct {
	Total       int64   `json:"total"`
	PaidRevenue float64 `json:"paid_revenue"`
	Pending     int64   `json:"pending"`
	Confirmed   int64   `json:"confirmed"`
	Processing  int64   `json:"processing"`
	Shipping    int64   `json:"shipping"`
	Completed   int64   `json:"completed"`
	Cancelled   int64   `json:"cancelled"`
	Refunded    int64   `json:"refunded"`
}

// CouponStats is the coupon summary.
type CouponStats struct {
	Total  int64 `json:"total"`
	Active int64 `json:"active"`
}

// UserStats is the customer summary.
type UserStats struct {
	Total int64 `json:"total"`
}

// StatsResponse is the dashboard overview payload.
type StatsResponse struct {
	Products ProductStats `json:"products"`
	Orders   OrderStats   `json:"orders"`
	Coupons  CouponStats  `json:"coupons"`
	Users    UserStats    `json:"users"`
}
