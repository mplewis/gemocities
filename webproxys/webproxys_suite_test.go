package webproxys_test

import (
	_ "embed"
	"testing"

	"github.com/mplewis/gemocities/webproxys"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

//go:embed test_in.gmi
var testIn string

//go:embed test_out.html
var testOut string

func TestWebproxys(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WebProxyS Suite")
}

var _ = Describe("ConvertToHTML", func() {
	It("converts as expected", func() {
		Expect(webproxys.ConvertToHTML(testIn)).To(Equal(testOut))
	})
})
