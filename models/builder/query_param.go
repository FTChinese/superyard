package builder

import (
	"fmt"
	"net/http"
	"strings"
)

// QueryParam contains a query parameter key-value pair.
type QueryParam struct {
	Name  string
	Value string
}

func NewQueryParam(key string) *QueryParam {
	return &QueryParam{Name: key}
}

func (p *QueryParam) SetValue(req *http.Request) *QueryParam {
	p.Value = req.Form.Get(p.Name)
	return p
}

func (p *QueryParam) Sanitize() *QueryParam {
	p.Name = strings.TrimSpace(p.Name)
	p.Value = strings.TrimSpace(p.Value)

	return p
}

func (p *QueryParam) Validate() error {
	if p.Value == "" {
		return fmt.Errorf("query parameter %s should have avlue", p.Name)
	}

	return nil
}
