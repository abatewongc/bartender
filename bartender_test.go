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

// TODO: test selectRandomChampionSkin

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

// TestSelectRandomSkin tests selectRandomChampionSkin works
func TestSelectRandomSkin(t *testing.T) {
	//TODO
}

// TestURLError tests selectRandomChampionSkin properly returns an error
// if the lcu client returns an error when creating a url.URL object
func TestURLError(t *testing.T) {
	//TODO
}

// TestGetSkinCarouselError tests selectRandomChampionSkin properly returns an
// error if there is an error getting the skin carousel
func TestGetSkinCarouselError(t *testing.T) {
	//TODO
}

// TestGetSkinCarouselBadResponse tests selectRandomChampionSkin properly
// returns an error on a non-200 status code when getting the skin
// carousel
func TestGetSkinCarouselBadResponse(t *testing.T) {
	//TODO
}

// TestEmptySkinList tests selectRandomChampionSkin properly returns an error
// if an empty skin list is returned from getting the skin carousel
func TestEmptySkinList(t *testing.T) {
	//TODO
}

// TestSelectSkinResponseError tests selectRandomChampionSkin properly returns
// an error if there is an error selecting the skin.
func TestSelectSkinResponseError(t *testing.T) {
	//TODO
}

// TestBadSelectSkinResponse tests selectRandomChampionSkin properly returns
// an error if a non-200 status code is returned while selecting a skin.
func TestBadSelectSkinResponse(t *testing.T) {
	//TODO
}
