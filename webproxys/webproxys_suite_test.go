package webproxys_test

import (
	_ "embed"
	"fmt"
	"strings"
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
		actual := strings.TrimSpace(webproxys.ConvertToHTML(testIn))
		expected := strings.TrimSpace(testOut)
		fmt.Println("ACTUAL:")
		fmt.Println("----------------------------------------")
		fmt.Println(actual)
		fmt.Println("----------------------------------------")
		fmt.Println("EXPECTED:")
		fmt.Println("----------------------------------------")
		fmt.Println(expected)
		fmt.Println("----------------------------------------")
		Expect(actual).To(Equal(expected))
	})
})
