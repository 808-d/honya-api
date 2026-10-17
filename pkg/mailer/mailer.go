package mailer

import (
	"context"
	"log"
)

// Mailer delivers transactional emails.
type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

// Dummy is a development Mailer that logs messages instead of sending them.
type Dummy struct{}

func NewDummy() *Dummy {
	return &Dummy{}
}

func (d *Dummy) Send(ctx context.Context, to, subject, body string) error {
	log.Printf("[dummy-email] to=%s subject=%q body=%s", to, subject, body)
	return nil
}
