package reader

import (
	"errors"
	"github.com/guregu/null"
	"strings"
)

type FtcInfo struct {
	ID    string `json:"id" db:"ftc_id"`
	Email string `json:"email" db:"email"`
}

type WxInfo struct {
	UnionID  string      `json:"unionId" db:"union_id"`
	Nickname null.String `json:"nickname" db:"nickname"`
}

type SearchParam struct {
	Email string `schema:"email"`
	Q     string `schema:"q"`
}

func (p *SearchParam) Sanitize() {
	p.Email = strings.TrimSpace(p.Email)
	p.Q = strings.TrimSpace(p.Q)
}

func (p SearchParam) RequireEmail() error {
	if p.Email == "" {
		return errors.New("email is required")
	}

	return nil
}

func (p SearchParam) RequireQ() error {
	if p.Q == "" {
		return errors.New("missing query parameter 'q'")
	}

	return nil
}
