package gemocities_test

import (
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
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gemocities Suite")
}

var _ = Describe("server", func() {
	var contentDir string
	var rq Requestor

	BeforeEach(func() {
		zerolog.SetGlobalLevel(zerolog.WarnLevel) // HACK

		cd, err := os.MkdirTemp("", "")
		Expect(err).ToNot(HaveOccurred())
		contentDir = cd
		gemSrv, err := geminis.BuildServer(geminis.ServerArgs{
			GeminiCertsDir: "test/certs",
			UserManager:    &user.Manager{Store: ez3.NewMemory(), TestMode: true},
			ContentManager: &content.Manager{Dir: contentDir},
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
		Expect(resp.Body()).To(ContainSubstring("This is the home page"))
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

	It("confirms registration details and creates an account with a parking page", func() {
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
	})
})
