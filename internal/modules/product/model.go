package product

import "time"

// Product is a catalog product (a book, manga or light novel).
type Product struct {
	ID          int64
	PublicID    string
	Title       string
	Cover       string
	Description string
	PublisherID *string
	PublishTime *time.Time
	Stock       int32
	Price       float64
	Status      int16
	Type        int16
	CreateTime  time.Time
	UpdateTime  time.Time
	AuthorIDs   []string
	GenreIDs    []string
	TagIDs      []string
}

// ProductWrite is a validated product payload ready to persist.
type ProductWrite struct {
	Title       string
	Cover       string
	Description string
	PublisherID *string
	PublishTime *time.Time
	Stock       int32
	Price       float64
	Status      int16
	Type        int16
	AuthorIDs   []string
	GenreIDs    []string
	TagIDs      []string
}
