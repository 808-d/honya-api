package tag

import "time"

// Tag is a catalog tag.
type Tag struct {
	ID         int64
	PublicID   string
	Name       string
	CreateTime time.Time
	UpdateTime time.Time
}
