package mail_test

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	"github.com/mplewis/gemocities/mail"
	"github.com/mplewis/gemocities/types"
	"github.com/mplewis/gemocities/user"
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
	"reset-password": "Reset your password: /reset-password?token={{ .Token }}",
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

	Describe("Send", func() {

		mailer := mail.Mailer{
			AppDomain: "amaya.com",
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
						Body:     "Reset your password: /reset-password?token=3458762345978",
					},
				},
			}))
		})
	})

	Describe("Mailer template methods", func() {
		mailer := mail.New(types.Config{
			AppDomain:    "amaya.com",
			SMTPUsername: "postmaster",
			SMTPPassword: "qubits",
			SMTPHost:     "mail.amaya.com",
			SMTPPort:     487,
		})

		Describe("SendVerificationEmail", func() {
			It("sends the expected email", func() {
				err := mailer.SendVerificationEmail(user.User{Email: "lily@amaya.com", Name: "lily", VerificationToken: "deadbeefcafe"})
				Expect(err).ToNot(HaveOccurred())
				Expect(sentMails).To(Equal([]SentMail{{
					SMTPArgs: mail.SMTPArgs{Host: "mail.amaya.com", Port: 487, Username: "postmaster", Password: "qubits"},
					Rendered: mail.Rendered{
						Headers: map[string][]string{
							"From":    []string{"Gemocities <welcome@amaya.com>"},
							"To":      []string{"lily@amaya.com"},
							"Subject": []string{"Confirm your Gemocities account"},
						},
						MimeType: "text/plain",
						Body: heredoc.Doc(`
							Hello! Please follow this link to verify your email address for your new Gemocities account ~lily:

							gemini://gemocities.com/account/verify?token=deadbeefcafe

							If you did not sign up for this account, you can safely ignore this email and the account will be deleted automatically.
						`),
					},
				}}))
			})
		})
	})
})
