package product

// ProductInput is the body of the create and update product endpoints.
type ProductInput struct {
	Title       string   `json:"title"`
	Cover       string   `json:"cover"`
	Description string   `json:"description"`
	PublisherID *string  `json:"publisher_id"`
	PublishTime *string  `json:"publish_time"`
	Stock       int32    `json:"stock"`
	Price       float64  `json:"price"`
	Status      int16    `json:"status"`
	Type        int16    `json:"type"`
	AuthorIDs   []string `json:"author_ids"`
	GenreIDs    []string `json:"genre_ids"`
	TagIDs      []string `json:"tag_ids"`
}

// ProductResponse is the public representation of a product.
type ProductResponse struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Cover       string   `json:"cover"`
	Description string   `json:"description"`
	PublisherID *string  `json:"publisher_id,omitempty"`
	PublishTime *string  `json:"publish_time,omitempty"`
	Stock       int32    `json:"stock"`
	Price       float64  `json:"price"`
	Status      int16    `json:"status"`
	Type        int16    `json:"type"`
	AuthorIDs   []string `json:"author_ids"`
	GenreIDs    []string `json:"genre_ids"`
	TagIDs      []string `json:"tag_ids"`
	CreateTime  string   `json:"create_time"`
	UpdateTime  string   `json:"update_time"`
}

// ListFilter holds the validated query parameters for listing products.
type ListFilter struct {
	Status   *int16
	Type     *int16
	Search   string
	MinPrice *float64
	MaxPrice *float64
	Page     int
	PerPage  int
}

// ListResult is a page of products.
type ListResult struct {
	Items   []ProductResponse `json:"items"`
	Total   int64             `json:"total"`
	Page    int               `json:"page"`
	PerPage int               `json:"per_page"`
}
