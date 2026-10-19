// Package mailer sends transactional email. It uses Gmail over SMTP with an
// app password by default, or logs messages when no credentials are configured.
package mailer

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"
)

// Message is an email carrying both a plain-text and an HTML body.
type Message struct {
	To      string
	Subject string
	Text    string
	HTML    string
}

// Mailer delivers transactional emails.
type Mailer interface {
	Send(ctx context.Context, msg Message) error
}

// Dummy is a development Mailer that logs messages instead of sending them.
type Dummy struct{}

func NewDummy() *Dummy {
	return &Dummy{}
}

func (d *Dummy) Send(ctx context.Context, msg Message) error {
	log.Printf("[dummy-email] to=%s subject=%q\n%s", msg.To, msg.Subject, msg.Text)
	return nil
}

// Gmail SMTP settings. These are fixed; only the email and app password vary.
const (
	smtpHost    = "smtp.gmail.com"
	smtpPort    = "587"
	smtpAddr    = smtpHost + ":" + smtpPort
	fromName    = "Honya"
	sendTimeout = 10 * time.Second
)

// SMTP sends mail through Gmail using an app password.
type SMTP struct {
	from string
	user string
	pass string
}

// NewSMTP builds a Gmail mailer from an email address and an app password.
func NewSMTP(email, password string) *SMTP {
	return &SMTP{
		from: email,
		user: email,
		pass: strings.ReplaceAll(password, " ", ""),
	}
}

// New returns a Gmail mailer when both credentials are present, otherwise a
// Dummy mailer that only logs.
func New(email, password string) Mailer {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return NewDummy()
	}
	return NewSMTP(email, password)
}

// Send renders the message and delivers it over an authenticated STARTTLS
// connection, honoring the context deadline.
func (s *SMTP) Send(ctx context.Context, msg Message) error {
	if _, err := mail.ParseAddress(msg.To); err != nil {
		return fmt.Errorf("invalid recipient address: %w", err)
	}
	if strings.ContainsAny(msg.To+msg.Subject, "\r\n") {
		return fmt.Errorf("invalid message header")
	}

	raw, err := buildMessage(s.from, msg)
	if err != nil {
		return err
	}
	return s.deliver(ctx, msg.To, raw)
}

func (s *SMTP) deliver(ctx context.Context, to string, raw []byte) error {
	ctx, cancel := context.WithTimeout(ctx, sendTimeout)
	defer cancel()

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", smtpAddr)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	if ok, _ := client.Extension("STARTTLS"); !ok {
		return fmt.Errorf("smtp server does not support STARTTLS")
	}
	if err := client.StartTLS(&tls.Config{ServerName: smtpHost}); err != nil {
		return err
	}

	if err := client.Auth(smtp.PlainAuth("", s.user, s.pass, smtpHost)); err != nil {
		return err
	}
	if err := client.Mail(s.from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(raw); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return client.Quit()
}

func buildMessage(from string, msg Message) ([]byte, error) {
	recipient, err := mail.ParseAddress(msg.To)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient address: %w", err)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	textPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/plain; charset=\"utf-8\""},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	if err != nil {
		return nil, err
	}
	if err := writeQuotedPrintable(textPart, msg.Text); err != nil {
		return nil, err
	}

	htmlPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/html; charset=\"utf-8\""},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	if err != nil {
		return nil, err
	}
	if err := writeQuotedPrintable(htmlPart, msg.HTML); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	var out bytes.Buffer
	writeHeader(&out, "From", fmt.Sprintf("%s <%s>", mime.QEncoding.Encode("utf-8", fromName), from))
	writeHeader(&out, "To", recipient.String())
	writeHeader(&out, "Subject", mime.QEncoding.Encode("utf-8", msg.Subject))
	writeHeader(&out, "Date", time.Now().Format(time.RFC1123Z))
	writeHeader(&out, "Message-ID", newMessageID(from))
	writeHeader(&out, "MIME-Version", "1.0")
	writeHeader(&out, "Content-Type", "multipart/alternative; boundary=\""+writer.Boundary()+"\"")
	out.WriteString("\r\n")
	out.Write(body.Bytes())

	return out.Bytes(), nil
}

func writeQuotedPrintable(w interface{ Write([]byte) (int, error) }, content string) error {
	qp := quotedprintable.NewWriter(w)
	if _, err := qp.Write([]byte(toCRLF(content))); err != nil {
		return err
	}
	return qp.Close()
}

func writeHeader(buf *bytes.Buffer, key, value string) {
	buf.WriteString(key)
	buf.WriteString(": ")
	buf.WriteString(value)
	buf.WriteString("\r\n")
}

// toCRLF normalizes line endings to CRLF as required by SMTP.
func toCRLF(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\n", "\r\n")
}

func newMessageID(from string) string {
	random := make([]byte, 16)
	_, _ = rand.Read(random)

	domain := "honya"
	if i := strings.LastIndex(from, "@"); i >= 0 {
		domain = from[i+1:]
	}
	return fmt.Sprintf("<%d.%s@%s>", time.Now().UnixNano(), hex.EncodeToString(random), domain)
}
