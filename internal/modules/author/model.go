package author

import "time"

// Author is a catalog author.
type Author struct {
	ID         int64
	PublicID   string
	Name       string
	CreateTime time.Time
	UpdateTime time.Time
}
