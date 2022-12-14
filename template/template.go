package template

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"text/template"
)

type Cache struct {
	*embed.FS
	Prefix    string
	Suffix    string
	templates map[string]*template.Template
}

func (c *Cache) load(name string) (*template.Template, error) {
	if c.templates == nil {
		c.templates = make(map[string]*template.Template)
	}
	if tmpl, ok := c.templates[name]; ok {
		return tmpl, nil
	}

	path := fmt.Sprintf("%s%s%s", c.Prefix, name, c.Suffix)
	tmpl, err := template.New(name).ParseFS(c.FS, path)
	if err != nil {
		return nil, fmt.Errorf("error parsing template at %s: %w", path, err)
	}
	c.templates[name] = tmpl
	return tmpl, nil
}

func (c *Cache) Render(dst io.Writer, name string, data any) error {
	tmpl, err := c.load(name)
	if err != nil {
		return fmt.Errorf("error loading template %s: %w", name, err)
	}
	err = tmpl.ExecuteTemplate(dst, name+c.Suffix, data)
	if err != nil {
		return fmt.Errorf("error executing template %s: %w", name, err)
	}
	return nil
}

func (c *Cache) RenderString(name string, data any) (string, error) {
	b := &bytes.Buffer{}
	err := c.Render(b, name, data)
	return b.String(), err
}
