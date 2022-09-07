package template_test

import (
	"embed"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/mplewis/gemocities/template"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

//go:embed test/templates/*
var templates embed.FS

func TestTemplate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Template Suite")
}

func getAllFilenames(fs *embed.FS, path string) (out []string, err error) {
	if len(path) == 0 {
		path = "."
	}
	entries, err := fs.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		fp := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			res, err := getAllFilenames(fs, fp)
			if err != nil {
				return nil, err
			}
			out = append(out, res...)
			continue
		}
		out = append(out, fp)
	}
	return
}

var _ = Describe("Cache", func() {
	It("works as intended", func() {
		fmt.Println(getAllFilenames(&templates, ""))

		c := &template.Cache{
			FS:     &templates,
			Prefix: "test/templates/",
			Suffix: ".md",
		}

		out, err := c.RenderString("homepage", nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(out).To(ContainSubstring("# Welcome!"))

		data := struct {
			Username string
			Balance  string
		}{"krina-alizond-114", "1,042.67"}
		out, err = c.RenderString("account", data)
		Expect(err).ToNot(HaveOccurred())
		Expect(out).To(ContainSubstring("You're logged in as krina-alizond-114"))
		Expect(out).To(ContainSubstring("Current balance: $1,042.67"))
	})
})
