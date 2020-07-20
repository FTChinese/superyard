package views

import (
	"errors"
	"github.com/labstack/echo/v4"
	"io"
)

type Views struct {
	templates map[string]string
}

func New() *Views {
	return &Views{templates: templates}
}

func (v *Views) Get(name string) (string, error) {
	tmpl, ok := templates[name]

	if !ok {
		return "", errors.New("html template not found")
	}

	return tmpl, nil
}

func (v *Views) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, err := v.Get(name)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(tmpl))
	if err != nil {
		return err
	}

	return nil
}
