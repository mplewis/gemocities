package mail

import (
	"gopkg.in/gomail.v2"
)

// templateCache describes a template.Cache.
type templateCache interface {
	RenderString(name string, data any) (string, error)
}

// Send is an interface for gomail that can be stubbed in testing.
var Send = func(s SMTPArgs, r Rendered) error {
	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)
	msg := gomail.NewMessage()
	for k, v := range r.Headers {
		msg.SetHeader(k, v...)
	}
	msg.SetBody(r.MimeType, r.Body)
	return d.DialAndSend(msg)
}

// Mailer sends emails using a library of templates.
type Mailer struct {
	SMTPArgs
	Templates templateCache
}

// Send sends an email with the given content.
func (m *Mailer) Send(args MailArgs) error {
	content, err := args.render(m.Templates)
	if err != nil {
		return err
	}
	return Send(m.SMTPArgs, content)
}

// SMTPArgs is the configuration for connecting to an SMTP server.
type SMTPArgs struct {
	Host     string
	Port     int
	Username string
	Password string
}

// MailArgs is the pre-rendered content of an email using a named template.
type MailArgs struct {
	From     string
	To       []string
	Subject  string
	Template string
	Data     any
}

// render renders a MailArgs into a Rendered mail which is ready to send.
func (m *MailArgs) render(tc templateCache) (Rendered, error) {
	body, err := tc.RenderString(m.Template, m.Data)
	if err != nil {
		return Rendered{}, err
	}
	headers := map[string][]string{"From": []string{m.From}, "To": m.To, "Subject": []string{m.Subject}}
	return Rendered{Headers: headers, MimeType: "text/plain", Body: body}, nil
}

// Rendered is the rendered headers and body content of an email.
type Rendered struct {
	Headers  map[string][]string
	MimeType string
	Body     string
}
