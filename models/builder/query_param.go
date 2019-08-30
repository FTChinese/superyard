package builder

import (
	"errors"
	"strings"
)

type SearchParam struct {
	Email string `schema:"email"`
	Name  string `schema:"name"`
	Q     string `schema:"q"`
}

func (p *SearchParam) Sanitize() {
	p.Email = strings.TrimSpace(p.Email)
	p.Name = strings.TrimSpace(p.Name)
	p.Q = strings.TrimSpace(p.Q)
}

func (p SearchParam) RequireEmail() error {
	if p.Email == "" {
		return errors.New("email is required")
	}

	return nil
}

func (p SearchParam) NameOrEmail() error {
	if p.Email == "" && p.Name == "" {
		return errors.New("email or name should be specified")
	}

	return nil
}

func (p SearchParam) RequireQ() error {
	if p.Q == "" {
		return errors.New("missing query parameter 'q'")
	}

	return nil
}
