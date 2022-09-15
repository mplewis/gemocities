package mail

import (
	"embed"

	"github.com/mplewis/gemocities/template"
	"github.com/mplewis/gemocities/types"
	"github.com/rs/zerolog/log"

	"gopkg.in/gomail.v2"
)

//go:embed templates/*
var templates embed.FS

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
	AppDomain string
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

func New(cfg types.Config) *Mailer {
	sa := SMTPArgs{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
	}
	tc := template.Cache{FS: &templates, Prefix: "templates/", Suffix: ".txt"}
	log.Info().Str("host", sa.Host).Int("port", sa.Port).Str("username", sa.Username).Msg("Mailer initialized")
	return &Mailer{SMTPArgs: sa, Templates: &tc, AppDomain: cfg.AppDomain}
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
