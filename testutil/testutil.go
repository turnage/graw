package testutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
)

// NewProxyClient returns an http.Client that redirects all requests to the
// redirect url.
func NewProxyClient(redirect string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return url.Parse(redirect)
			},
		},
	}
}

// NewServerFromResponse returns an httptest.Server that always responds with
// response and status.
func NewServerFromResponse(status int, response []byte) *httptest.Server {
	responseString := bytes.NewBuffer(response).String()
	responseWriter := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, responseString)
	}
	return httptest.NewServer(http.HandlerFunc(responseWriter))
}

// RepsonseIs returns true iff the response status code and body are identical
// to the provided expectations.
func ResponseIs(resp *http.Response, status int, expected []byte) bool {
	if resp.StatusCode != status {
		return false
	}

	if resp.Body == nil {
		if expected == nil {
			return true
		}
		return false
	}

	defer resp.Body.Close()
	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(actual, expected)
}
