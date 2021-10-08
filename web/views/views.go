package views

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type Views struct {
	templates map[string]string
}

func New() *Views {
	return &Views{templates: templates}
}

func (v *Views) Get(name string) (*template.Template, error) {
	tmpl, ok := templates[name]

	if !ok {
		return nil, fmt.Errorf("template [%s] not found", name)
	}

	return template.New(name).Parse(tmpl)
}

func (v *Views) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, err := v.Get(name)
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		return err
	}

	return nil
}
