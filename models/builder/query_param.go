package builder

import (
	"fmt"
	"net/http"
	"strings"
)

// QueryParam contains a query parameter key-value pair.
// This is used to handle cases like the a pair is allowed to have different keys, but always has a single value.
// The first key found with a value in the request form will
// be used. All other keys are ignored.
//
// Example
// the request query could be one of the two
// ?email=name@example.org
// ?name=user_name
type QueryParam struct {
	keys  []string
	Name  string
	Value string
}

// NewQueryParam creates a new instance.
func NewQueryParam(keys ...string) *QueryParam {
	p := QueryParam{}
	for _, v := range keys {
		p.keys = append(p.keys, v)
	}

	return &p
}

// SetValues tries to find a value for any of the specified keys.
func (p *QueryParam) SetValue(req *http.Request) *QueryParam {
	for _, k := range p.keys {
		v := req.Form.Get(k)
		if v != "" {
			p.Name = k
			p.Value = v
			break
		}
	}

	return p
}

// Sanitize removes empty strings.
func (p *QueryParam) Sanitize() *QueryParam {
	p.Name = strings.TrimSpace(p.Name)
	p.Value = strings.TrimSpace(p.Value)

	return p
}

// Validate ensures the Value is not empty.
func (p *QueryParam) Validate() error {
	if p.Value == "" {
		return fmt.Errorf("query parameter %s should have a value", p.Name)
	}

	return nil
}
