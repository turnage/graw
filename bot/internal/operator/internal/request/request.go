package request

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

// New returns an http.Request with the method, url, and values specified.
func New(method, url string, vals *url.Values) (*http.Request, error) {
	if method == "GET" {
		return buildGetRequest(url, vals)
	} else if method == "POST" {
		return buildPostRequest(url, vals)
	} else {
		return nil, fmt.Errorf("unsupported request method")
	}
}

// buildGetRequest returns an http.Request with the given url and GET form
// values set.
func buildGetRequest(url string, vals *url.Values) (*http.Request, error) {
	reqURL := url
	if vals != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, vals.Encode())
	}
	return http.NewRequest("GET", reqURL, nil)
}

// buildPostRequest returns an http.Request with the given url and POST form
// values set.
func buildPostRequest(url string, vals *url.Values) (*http.Request, error) {
	if vals == nil {
		return nil, fmt.Errorf("no values for POST body")
	}

	reqBody := bytes.NewBufferString(vals.Encode())
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	return req, nil
}
