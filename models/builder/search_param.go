package builder

import (
	"errors"
	"net/url"
	"strings"
)

// SearchParams defines the keys to search a staff in query parameters.
type SearchParams struct {
	Key   string
	Value string
}

func (p *SearchParams) Validate() error {
	if p.Key == "" {
		return errors.New("no search key present")
	}
	if p.Value == "" {
		return errors.New("no search value specified")
	}

	return nil
}

func NewSearchParam(params url.Values, keys []string) SearchParams {
	p := SearchParams{}

	for _, k := range keys {
		if v := strings.TrimSpace(params.Get(k)); v != "" {
			p.Key = k
			p.Value = v
			break
		}
	}

	return p
}
