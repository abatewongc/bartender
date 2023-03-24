package bartender

import (
	"bartender/internal/mock"
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errEmpty = errors.New("")

// Test isChampionLocked properly returns true
func TestIsChampionLocked(t *testing.T) {
	client := mock.Client{
		URLResponse: url.URL{},
		URLError:    nil,
		GetResponse: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
		},
		GetError: nil,
	}
	svc := New(client)

	actual, err := svc.isChampionLocked()

	assert.NoError(t, err)
	assert.True(t, actual)
}

// Test isChampionLocked returns error on error fetching URL
func TestIsChampionLockedErrorURLError(t *testing.T) {
	client := mock.Client{
		URLResponse: url.URL{},
		URLError:    errEmpty,
	}
	svc := New(client)

	_, err := svc.isChampionLocked()

	assert.Error(t, err)
}

// Test isChampionLocked returns error on error sending get request
func TestIsChampionLockedGetError(t *testing.T) {
	client := mock.Client{
		GetResponse: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
		},
		GetError: errEmpty,
	}
	svc := New(client)

	_, err := svc.isChampionLocked()

	assert.Error(t, err)
}

// Test isChampionLocked returns error on a bad status code from the get request
func TestIsChampionLockedGetBadStatusError(t *testing.T) {
	client := mock.Client{
		GetResponse: &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
		},
		GetError: nil,
	}
	svc := New(client)

	_, err := svc.isChampionLocked()

	assert.Error(t, err)
}
