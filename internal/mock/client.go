package mock

import (
	"net/http"
	"net/url"
)

// Client mocks the client.Client interface
type Client struct {
	Called             map[string]int
	NewRequestResponse *http.Request
	NewRequestError    error
	URLResponse        url.URL
	URLError           error
	GetResponse        *http.Response
	GetError           error
	PostResponse       *http.Response
	PostError          error
}

func (c Client) NewRequest(s string, u url.URL, b []byte) (*http.Request, error) {
	c.Called["NewRequest"]++
	return c.NewRequestResponse, c.NewRequestError
}

func (c Client) URL(uri string) (url.URL, error) {
	c.Called["URL"]++
	return c.URLResponse, c.URLError
}

func (c Client) Get(url.URL) (*http.Response, error) {
	c.Called["Get"]++
	return c.GetResponse, c.GetError
}

func (c Client) Post(url.URL, []byte) (*http.Response, error) {
	c.Called["Post"]++
	return c.PostResponse, c.PostError
}
