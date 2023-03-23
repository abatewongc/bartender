package internal

import (
	"net/http"
	"net/url"
)

type MockClient struct{}

func (mc MockClient) NewRequest(s string, u url.URL, b []byte) (*http.Request, error) {
	return &http.Request{}, nil
}

func (mc MockClient) URL(uri string) (url.URL, error) {
	return url.URL{}, nil
}

func (mc MockClient) Get(url.URL) (*http.Response, error) {
	return &http.Response{}, nil
}

func (mc MockClient) Post(url.URL, []byte) (*http.Response, error) {
	return &http.Response{}, nil
}
