package mail

import (
	"embed"
	"fmt"

	"github.com/mplewis/gemocities/user"

	"gopkg.in/gomail.v2"
)

//go:embed templates/*
var templates embed.FS

// IMailer is a high-level interface for sending domain mail.
type IMailer interface {
	SendVerificationEmail(user user.User) error
}

// Send is an interface for gomail that can be stubbed in testing.
var Send = func(s SMTPArgs, r Rendered) error {
	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)
	msg := gomail.NewMessage()
	for k, v := range r.Headers {
		msg.SetHeader(k, v...)
	}
	msg.SetBody(r.MimeType, r.Body)
	err := d.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("error sending mail: %w", err)
	}
	return nil
}
