package gemocities_test

import (
	"context"
	"net/url"
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

var _ = Describe("Gemocities", func() {
	It("passes a basic request test", func() {
		ctx := context.Background()
		contentDir, err := os.MkdirTemp("", "")
		Expect(err).ToNot(HaveOccurred())
		defer os.RemoveAll(contentDir)

		umgr := &user.Manager{Store: ez3.NewMemory()}
		cmgr := &content.Manager{Dir: contentDir}

		gemSrv, err := geminis.BuildServer(geminis.ServerArgs{
			GeminiCertsDir: "test/certs",
			GeminiHost:     ":1965",
			UserManager:    umgr,
			ContentManager: cmgr,
			ContentDir:     contentDir,
		})
		Expect(err).ToNot(HaveOccurred())

		u, err := url.Parse("/")
		Expect(err).ToNot(HaveOccurred())
		req := gemini.Request{URL: u}
		var resp ResponseBuffer
		gemSrv.Handler.ServeGemini(ctx, &resp, &req)
		Expect(resp.Status).To(Equal(gemini.StatusSuccess))
		Expect(resp.Body()).To(ContainSubstring("This is the home page"))
	})
})
