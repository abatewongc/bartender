package mock

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

// client mocks the client.client and bartender.HTTPclient interfaces
type client struct {
	NewRequestResponse *http.Request
	NewRequestError    error
	URLResponse        url.URL
	URLError           error
	GetResponse        *http.Response
	GetError           error
	PostResponse       *http.Response
	PostError          error
	CalledMethods      map[string]int
	DoRequestRecorder  *http.Request
}

func Newclient() *client {
	return &client{
		NewRequestResponse: &http.Request{},
		NewRequestError:    nil,
		URLResponse:        url.URL{},
		URLError:           nil,
		GetResponse:        &http.Response{Body: io.NopCloser(bytes.NewBuffer([]byte{}))},
		GetError:           nil,
		PostResponse:       &http.Response{Body: io.NopCloser(bytes.NewBuffer([]byte{}))},
		PostError:          nil,
		CalledMethods:      map[string]int{},
		DoRequestRecorder:  &http.Request{},
	}
}

func (c client) NewRequest(s string, u url.URL, b []byte) (*http.Request, error) {
	c.CalledMethods["NewRequest"]++
	return c.NewRequestResponse, c.NewRequestError
}

func (c client) URL(uri string) (url.URL, error) {
	c.CalledMethods["URL"]++
	return c.URLResponse, c.URLError
}

func (c client) Get(url.URL) (*http.Response, error) {
	c.CalledMethods["Get"]++
	return c.GetResponse, c.GetError
}

func (c client) Post(url.URL, []byte) (*http.Response, error) {
	c.CalledMethods["Post"]++
	return c.PostResponse, c.PostError
}

func (c client) Do(req *http.Request) (*http.Response, error) {
	c.CalledMethods["Do"]++
	switch req.Method {
	case "Patch":
		c.CalledMethods["Patch"]++
	}
	return &http.Response{}, nil
}
