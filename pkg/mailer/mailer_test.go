package mailer

import (
	"strings"
	"testing"
)

func TestBuildMessage(t *testing.T) {
	raw, err := buildMessage("no-reply@honya.test", Message{
		To:      "user@example.com",
		Subject: "Verify your Honya account",
		Text:    "Hi,\n\nOpen: https://honya.test/verify?token=abc\n",
		HTML:    "<p>Hi,</p>\n<p><a href=\"https://honya.test/verify?token=abc\">Verify</a></p>\n",
	})
	if err != nil {
		t.Fatalf("buildMessage: %v", err)
	}
	msg := string(raw)

	for _, want := range []string{
		"From: ",
		"no-reply@honya.test",
		"To: <user@example.com>",
		"Subject: ",
		"MIME-Version: 1.0",
		"multipart/alternative; boundary=",
		`Content-Type: text/plain; charset="utf-8"`,
		`Content-Type: text/html; charset="utf-8"`,
		"Content-Transfer-Encoding: quoted-printable",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("message missing %q", want)
		}
	}

	if bare := strings.ReplaceAll(msg, "\r\n", ""); strings.Contains(bare, "\n") {
		t.Error("message contains bare LF; want CRLF only")
	}
}

func TestNewSelection(t *testing.T) {
	if _, ok := New("", "").(*Dummy); !ok {
		t.Error("want Dummy when credentials are missing")
	}
	if _, ok := New("user@gmail.com", "app password").(*SMTP); !ok {
		t.Error("want SMTP when credentials are present")
	}
}

func TestNewSMTPStripsSpaces(t *testing.T) {
	s := NewSMTP("user@gmail.com", "abcd efgh ijkl mnop")
	if s.pass != "abcdefghijklmnop" {
		t.Errorf("password = %q, want spaces stripped", s.pass)
	}
	if s.user != "user@gmail.com" || s.from != "user@gmail.com" {
		t.Errorf("user/from = %q/%q", s.user, s.from)
	}
}
