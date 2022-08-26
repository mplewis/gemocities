package geminis

import (
	"embed"
	"fmt"
	"io"
	"path"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

type TemplateCache struct {
	templates map[string]*template.Template
}

func (t *TemplateCache) load(name string) (*template.Template, error) {
	if t.templates == nil {
		t.templates = make(map[string]*template.Template)
	}
	if tmpl, ok := t.templates[name]; ok {
		return tmpl, nil
	}

	p := path.Join("templates", fmt.Sprintf("%s.gmi", name))
	raw, err := templateFS.ReadFile(p)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(name).Parse(string(raw))
	if err != nil {
		return nil, err
	}
	t.templates[name] = tmpl
	return tmpl, nil
}

func (t *TemplateCache) Render(dst io.Writer, name string, data any) error {
	tmpl, err := t.load(name)
	if err != nil {
		return err
	}
	return tmpl.Execute(dst, data)
}
