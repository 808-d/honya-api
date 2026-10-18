package publisher

import "time"

// Publisher is a catalog publisher.
type Publisher struct {
	ID         int64
	PublicID   string
	Name       string
	CreateTime time.Time
	UpdateTime time.Time
}
