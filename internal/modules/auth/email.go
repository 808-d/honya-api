package auth

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	texttemplate "text/template"
	"time"

	"github.com/808-d/honya-api/pkg/mailer"
)

// emailData is the data passed to the shared email templates.
type emailData struct {
	Username string
	Action   string
	Link     string
	Expires  string
}

const emailTextTmpl = `Hi {{.Username}},

Please {{.Action}} by opening this link:

{{.Link}}

This link will expire in {{.Expires}}. If you did not request this, you can safely ignore this email.

- Honya
`

const emailHTMLTmpl = `<!DOCTYPE html>
<html>
  <body style="font-family: Arial, Helvetica, sans-serif; color: #222222; line-height: 1.5;">
    <p>Hi {{.Username}},</p>
    <p>Please {{.Action}} by clicking the button below:</p>
    <p>
      <a href="{{.Link}}"
         style="display: inline-block; padding: 10px 18px; background: #111111; color: #ffffff; text-decoration: none; border-radius: 6px;">
        {{.Action}}
      </a>
    </p>
    <p>Or copy this link into your browser:<br><a href="{{.Link}}">{{.Link}}</a></p>
    <p>This link will expire in {{.Expires}}. If you did not request this, you can safely ignore this email.</p>
    <p>- Honya</p>
  </body>
</html>
`

var (
	emailText = texttemplate.Must(texttemplate.New("text").Parse(emailTextTmpl))
	emailHTML = htmltemplate.Must(htmltemplate.New("html").Parse(emailHTMLTmpl))
)

func buildEmail(to, subject, username, action, link, expires string) (mailer.Message, error) {
	data := emailData{Username: username, Action: action, Link: link, Expires: expires}

	var text, html bytes.Buffer
	if err := emailText.Execute(&text, data); err != nil {
		return mailer.Message{}, err
	}
	if err := emailHTML.Execute(&html, data); err != nil {
		return mailer.Message{}, err
	}

	return mailer.Message{
		To:      to,
		Subject: subject,
		Text:    text.String(),
		HTML:    html.String(),
	}, nil
}

// expiryText renders a token lifetime for humans, e.g. "30 minutes".
func expiryText(ttl time.Duration) string {
	minutes := int(ttl.Round(time.Minute).Minutes())
	if minutes == 1 {
		return "1 minute"
	}
	return fmt.Sprintf("%d minutes", minutes)
}

// verificationEmail is sent after registration to confirm the email address.
func verificationEmail(to, username, link, expires string) (mailer.Message, error) {
	return buildEmail(to, "Verify your Honya account", username, "confirm your email address", link, expires)
}

// passwordResetEmail is sent when a user requests a password reset.
func passwordResetEmail(to, username, link, expires string) (mailer.Message, error) {
	return buildEmail(to, "Reset your Honya password", username, "reset your password", link, expires)
}
