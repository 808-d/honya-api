package tag

// Input is the body of the create and update tag endpoints.
type Input struct {
	Name string `json:"name"`
}

// Response is the public representation of a tag.
type Response struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

// ListResult is a page of tags.
type ListResult struct {
	Items   []Response `json:"items"`
	Total   int64      `json:"total"`
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
}
