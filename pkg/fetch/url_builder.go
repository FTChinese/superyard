package fetch

import (
	"net/url"
	"strconv"
	"strings"
)

type URLBuilder struct {
	base  string
	paths []string
	query url.Values
}

func NewURLBuilder(base string) URLBuilder {
	return URLBuilder{
		base:  strings.TrimSuffix(base, "/"),
		paths: make([]string, 0),
		query: make(url.Values),
	}
}

func (b URLBuilder) AddPath(p string) URLBuilder {
	p = strings.TrimSuffix(strings.TrimPrefix(p, "/"), "/")

	b.paths = append(b.paths, p)
	return b
}

func (b URLBuilder) AddQuery(k, v string) URLBuilder {
	b.query.Add(k, v)
	return b
}

func (b URLBuilder) AddQueryBool(k string, v bool) URLBuilder {
	b.query.Add(k, strconv.FormatBool(v))
	return b
}

func (b URLBuilder) String() string {
	var buf strings.Builder
	if b.base != "" {
		buf.WriteString(b.base)
	}

	path := strings.Join(b.paths, "/")
	if path != "" {
		buf.WriteByte('/')
		buf.WriteString(path)
	}

	query := b.query.Encode()

	if query != "" {
		buf.WriteByte('?')
		buf.WriteString(query)
	}

	return buf.String()
}
