package letter

import (
	"fmt"
	"strings"
	"text/template"
)

var tmplCache = map[string]*template.Template{}

const (
	keySignUp         = "signUp"
	keyPwReset        = "passwordReset"
	keyOrderConfirmed = "orderConfirmed"
)

type CtxSignUp struct {
	DisplayName string
	LoginName   string
	Password    string
	LoginURL    string
}

type CtxPasswordReset struct {
	DisplayName string
	URL         string
}

type CtxConfirmOrder struct {
	Name           string
	OrderCreatedAt string
	OrderID        string
	OrderAmount    string
	PayMethod      string
	OrderStartDate string
	OrderEndDate   string
	Tier           string
	ExpirationDate string
}

func Render(name string, ctx interface{}) (string, error) {
	tmpl, ok := tmplCache[name]
	var err error
	if !ok {
		tmplStr, ok := templates[name]
		if !ok {
			return "", fmt.Errorf("template %s not found", name)
		}

		tmpl, err = template.New(name).Parse(tmplStr)
		if err != nil {
			return "", err
		}
		tmplCache[name] = tmpl
	}

	var body strings.Builder
	err = tmpl.Execute(&body, ctx)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}

func RenderSignUp(ctx CtxSignUp) (string, error) {
	return Render(keySignUp, ctx)
}

func RenderPasswordReset(ctx CtxPasswordReset) (string, error) {
	return Render(keyPwReset, ctx)
}

func RenderOrderConfirmed(ctx CtxConfirmOrder) (string, error) {
	return Render(keyOrderConfirmed, ctx)
}
