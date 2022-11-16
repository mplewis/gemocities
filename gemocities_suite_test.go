package gemocities_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/ez3"
	"github.com/mplewis/gemocities/content"
	"github.com/mplewis/gemocities/geminis"
	"github.com/mplewis/gemocities/user"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
)

func TestGemocities(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gemocities Suite")
}

type SentMail struct {
	To    string
	Token string
}

type FakeMailer struct {
	SentMails      []SentMail
	CrashWithError bool
}

var errSendEmailFailed = errors.New("sending email failed")

func (f *FakeMailer) SendVerificationEmail(user user.User) error {
	if f.CrashWithError {
		return errSendEmailFailed
	}
	f.SentMails = append(f.SentMails, SentMail{To: user.Email, Token: user.VerificationToken})
	return nil
}

var _ = Describe("server", func() {
	var contentDir string
	var rq Requestor
	var um user.Manager
	var cm content.Manager
	var fm FakeMailer

	BeforeEach(func() {
		zerolog.SetGlobalLevel(zerolog.WarnLevel) // HACK

		cd, err := os.MkdirTemp("", "")
		Expect(err).ToNot(HaveOccurred())
		contentDir = cd

		um = user.Manager{TestMode: true, Store: ez3.NewMemory()}
		cm = content.Manager{Dir: contentDir}
		fm = FakeMailer{}
		gemSrv, err := geminis.BuildServer(geminis.ServerArgs{
			GeminiCertsDir: "test/certs",
			UserManager:    &um,
			ContentManager: &cm,
			Mailer:         &fm,
			ContentDir:     contentDir,
		})
		Expect(err).ToNot(HaveOccurred())
		rq = Requestor{gemSrv}
	})

	AfterEach(func() {
		os.RemoveAll(contentDir)
	})

	It("presents the home page", func() {
		resp := rq.Request("/", nil)
		Expect(resp.Status).To(Equal(gemini.StatusSuccess))
		Expect(resp.Body()).To(ContainSubstring("# Welcome to Gemocities"))
	})

	It("requires certs for the account page", func() {
		resp := rq.Request("/account", nil)
		Expect(resp.Status).To(Equal(gemini.StatusCertificateRequired))
	})

	It("prompts new users to set up their account", func() {
		resp := rq.Request("/account", ClientCerts())
		Expect(resp.Status).To(Equal(gemini.StatusSuccess))
		Expect(resp.Body()).To(ContainSubstring("client certificate is not yet associated"))
	})

	It("requests registration details", func() {
		resp := rq.Request("/account/register", ClientCerts())
		Expect(resp.Status).To(Equal(gemini.StatusInput))
		Expect(resp.Meta).To(ContainSubstring("Enter your desired username"))
	})

	It("confirms registration details, creates an account with a parking page, and verifies", func() {
		resp := rq.RequestInput("/account/register", ClientCerts(), "elliot:mrr@fs0cie.ty")
		Expect(resp.Status).To(Equal(gemini.StatusSuccess))
		Expect(resp.Body()).To(ContainSubstring("confirm your new account details"))
		Expect(resp.Body()).To(ContainSubstring("Username: elliot"))
		Expect(resp.Body()).To(ContainSubstring("Email address: mrr@fs0cie.ty"))

		link, ok := resp.Links().WithText("Confirm and register")
		Expect(ok).To(BeTrue())
		Expect(link.URL).To(Equal("/account/register/confirm?username=elliot&email=mrr@fs0cie.ty"))

		resp = rq.Request(link.URL, ClientCerts())
		Expect(resp.Status).To(Equal(gemini.StatusRedirect))
		Expect(resp.Meta).To(Equal("/account"))

		resp = rq.Request("/~elliot/", ClientCerts())
		Expect(resp.Status).To(Equal(gemini.StatusSuccess))
		Expect(resp.Body()).To(ContainSubstring("This is your new user directory"))

		// user is unverified
		resp = rq.Request("/account", ClientCerts())
		Expect(resp.Status).To(Equal(gemini.StatusSuccess))
		Expect(resp.Body()).To(ContainSubstring("Your email address has not been verified"))

		// now verify the user
		Expect(fm.SentMails).To(HaveLen(1))
		mail := fm.SentMails[0]
		Expect(mail.To).To(Equal("mrr@fs0cie.ty"))

		resp = rq.Request(fmt.Sprintf("/account/verify?token=%s", mail.Token), ClientCerts())
		Expect(resp.Status).To(Equal(gemini.StatusRedirect))
		Expect(resp.Meta).To(Equal("/account"))

		resp = rq.Request("/account", ClientCerts())
		Expect(resp.Status).To(Equal(gemini.StatusSuccess))
		Expect(resp.Body()).To(ContainSubstring("Sign into WebDAV with the following credentials"))
	})

	Context("when sending the verification email fails", func() {
		BeforeEach(func() {
			fm.CrashWithError = true
		})
		AfterEach(func() {
			fm.CrashWithError = false
		})

		It("rolls back user account creation", func() {
			resp := rq.RequestInput("/account/register", ClientCerts(), "elliot:mrr@fs0cie.ty")
			Expect(resp.Status).To(Equal(gemini.StatusSuccess))
			Expect(resp.Body()).To(ContainSubstring("confirm your new account details"))
			Expect(resp.Body()).To(ContainSubstring("Username: elliot"))
			Expect(resp.Body()).To(ContainSubstring("Email address: mrr@fs0cie.ty"))

			link, ok := resp.Links().WithText("Confirm and register")
			Expect(ok).To(BeTrue())
			Expect(link.URL).To(Equal("/account/register/confirm?username=elliot&email=mrr@fs0cie.ty"))

			resp = rq.Request(link.URL, ClientCerts())
			Expect(resp.Status).To(Equal(gemini.StatusTemporaryFailure))
			Expect(resp.Meta).To(Equal("Sorry, there was an error creating your account. Please try again."))

			// verify rollback
			users, err := um.Store.List("")
			Expect(err).ToNot(HaveOccurred())
			Expect(users).To(BeEmpty())
			Expect(cm.Exists("elliot")).To(BeFalse())
		})
	})
})
