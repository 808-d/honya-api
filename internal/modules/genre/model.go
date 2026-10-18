package genre

import "time"

// Genre is a catalog genre.
type Genre struct {
	ID         int64
	PublicID   string
	Name       string
	CreateTime time.Time
	UpdateTime time.Time
}
