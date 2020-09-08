package fetch

import (
	"io"
	"net/http"
)

var httpClient = &http.Client{}

type RequestBuilder struct {
	method string
	url    string
	body   io.Reader
	Header http.Header
}

func NewRequest() *RequestBuilder {
	return &RequestBuilder{
		body: nil,
	}
}

func (b *RequestBuilder) Get(url string) *RequestBuilder {
	b.method = "GET"
	b.url = url

	return b
}

func (b *RequestBuilder) Post(url string) *RequestBuilder {
	b.method = "POST"
	b.url = url

	return b
}

func (b *RequestBuilder) Put(url string) *RequestBuilder {
	b.method = "PUT"
	b.url = url

	return b
}

func (b *RequestBuilder) SetAuth(key string) *RequestBuilder {
	b.Header.Add("Authorization", "Bearer "+key)

	return b
}

func (b *RequestBuilder) Send(body io.Reader) *RequestBuilder {
	b.body = body
	return b
}

func (b *RequestBuilder) End() (*http.Response, error) {
	req, err := http.NewRequest(b.method, b.url, b.body)
	if err != nil {
		return nil, err
	}

	req.Header = b.Header

	return httpClient.Do(req)
}
