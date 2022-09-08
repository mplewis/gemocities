package mail_test

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"github.com/mplewis/gemocities/mail"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMail(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mail Suite")
}

type SentMail struct {
	mail.SMTPArgs
	mail.Rendered
}

type fakeTemplateCache struct{}

var fakeTemplates = map[string]string{
	"reset-password": "Reset your password: /reset-password?user={{ .User }}&token={{ .Token }}",
}

func (f *fakeTemplateCache) RenderString(name string, data any) (string, error) {
	tbody, ok := fakeTemplates[name]
	if !ok {
		return "", fmt.Errorf("unknown template: %s", name)
	}
	b := bytes.Buffer{}
	err := template.Must(template.New("").Parse(tbody)).Execute(&b, data)
	return b.String(), err
}

var _ = Describe("Mailer", func() {
	var sentMails []SentMail
	BeforeEach(func() {
		sentMails = []SentMail{}
	})
	mail.Send = func(s mail.SMTPArgs, r mail.Rendered) error {
		sentMails = append(sentMails, SentMail{SMTPArgs: s, Rendered: r})
		return nil
	}

	mailer := mail.Mailer{
		SMTPArgs: mail.SMTPArgs{
			Host:     "mail.amaya.com",
			Port:     487,
			Username: "forest",
			Password: "oneworld",
		},
		Templates: &fakeTemplateCache{},
	}

	It("sends the expected mail", func() {
		err := mailer.Send(mail.MailArgs{
			From:     "welcome@amaya.com",
			To:       []string{"lily@amaya.com"},
			Subject:  "Reset your password",
			Template: "reset-password",
			Data: struct {
				User  string
				Token string
			}{User: "lily", Token: "3458762345978"},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(sentMails).To(Equal([]SentMail{
			{
				mail.SMTPArgs{
					Host:     "mail.amaya.com",
					Port:     487,
					Username: "forest",
					Password: "oneworld",
				},
				mail.Rendered{
					Headers: map[string][]string{
						"From":    []string{"welcome@amaya.com"},
						"To":      []string{"lily@amaya.com"},
						"Subject": []string{"Reset your password"},
					},
					MimeType: "text/plain",
					Body:     "Reset your password: /reset-password?user=lily&token=3458762345978",
				},
			},
		}))
	})
})
