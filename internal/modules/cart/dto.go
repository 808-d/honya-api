package cart

// AddItemRequest is the body of the add-to-cart endpoint.
type AddItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

// CartItemResponse is the public representation of a cart line. Every cart
// endpoint returns this same shape so the client can update its state directly.
type CartItemResponse struct {
	ID         int64  `json:"id"`
	Cover      string `json:"cover"`
	Title      string `json:"title"`
	Quantity   int32  `json:"quantity"`
	CreateTime string `json:"create_time"`
}

// CartResponse is the full cart.
type CartResponse struct {
	Items     []CartItemResponse `json:"items"`
	CartCount int64              `json:"cart_count"`
}

// CartMutationResponse is returned by add/remove so the client can patch the
// affected line and update the cart badge in one shot.
type CartMutationResponse struct {
	Item      CartItemResponse `json:"item"`
	CartCount int64            `json:"cart_count"`
}
