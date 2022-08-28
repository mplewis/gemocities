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
)

func TestGemocities(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gemocities Suite")
}

var _ = Describe("server", func() {
	var contentDir string
	var rq Requestor

	BeforeEach(func() {
		cd, err := os.MkdirTemp("", "")
		Expect(err).ToNot(HaveOccurred())
		contentDir = cd
		gemSrv, err := geminis.BuildServer(geminis.ServerArgs{
			GeminiCertsDir: "test/certs",
			UserManager:    &user.Manager{Store: ez3.NewMemory()},
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
})
