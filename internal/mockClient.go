package internal

import (
	"net/http"
	"net/url"
)

// MockClient mocks the client.Client interface
type MockClient struct {
	NewRequestResponse *http.Request
	NewRequestError    error
	URLResponse        url.URL
	URLError           error
	GetResponse        *http.Response
	GetError           error
	PostResponse       *http.Response
	PostError          error
}

func (mc MockClient) NewRequest(s string, u url.URL, b []byte) (*http.Request, error) {
	return mc.NewRequestResponse, mc.NewRequestError
}

func (mc MockClient) URL(uri string) (url.URL, error) {
	return mc.URLResponse, mc.URLError
}

func (mc MockClient) Get(url.URL) (*http.Response, error) {
	return mc.GetResponse, mc.GetError
}

func (mc MockClient) Post(url.URL, []byte) (*http.Response, error) {
	return mc.PostResponse, mc.PostError
}
