package cart

import "time"

// CartItem is a single product line in a user's cart.
type CartItem struct {
	ID         int64
	Cover      string
	Title      string
	Quantity   int32
	CreateTime time.Time
}
