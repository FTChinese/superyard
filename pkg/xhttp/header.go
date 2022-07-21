package xhttp

import "net/http"

func HeaderStaffName(n string) (string, string) {
	return XStaffName, n
}

func HeaderFtcID(id string) (string, string) {
	return XUserID, id
}

func HeaderWxID(id string) (string, string) {
	return XUnionID, id
}

type HeaderBuilder struct {
	h http.Header
}

func NewHeaderBuilder() *HeaderBuilder {
	return &HeaderBuilder{
		h: http.Header{},
	}
}

func (b *HeaderBuilder) WithFtcID(id string) *HeaderBuilder {
	b.h.Set(XUserID, id)
	return b
}

func (b *HeaderBuilder) WithUnionID(id string) *HeaderBuilder {
	b.h.Set(XUnionID, id)
	return b
}

func (b *HeaderBuilder) Build() http.Header {
	return b.h
}
