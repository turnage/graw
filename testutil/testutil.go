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

// TestRoundTripper redirects all requests to a test server.
type TestRoundTripper struct {
	URL *url.URL
}

func (t TestRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL = t.URL
	return http.DefaultClient.Do(r)
}

// NewProxyClient returns an http.Client that redirects all requests to the
// redirect url.
func NewProxyClient(redirURL *url.URL) *http.Client {
	return &http.Client{
		Transport: TestRoundTripper{URL: redirURL},
	}
}

// NewServerFromResponse returns an httptest.Server that always responds with
// response and status.
func NewServerFromResponse(stat int, resp []byte) (*httptest.Server, *url.URL) {
	responseString := bytes.NewBuffer(resp).String()
	responseWriter := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(stat)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, responseString)
	}
	server := httptest.NewServer(http.HandlerFunc(responseWriter))
	serverURL, _ := url.Parse(server.URL)
	return server, serverURL
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
