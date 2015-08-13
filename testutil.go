package graw

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
)

type bytesCloser struct {
	Buffer *bytes.Buffer
	Err    error
}

func (b bytesCloser) Read(p []byte) (int, error) {
	if b.Err != nil {
		return 0, b.Err
	}

	return b.Buffer.Read(p)
}

func (b bytesCloser) Close() error {
	return nil
}

// newServerFromResponse returns an httptest.Server that always responds with
// response and status.
func newServerFromResponse(stat int, resp []byte) *httptest.Server {
	responseString := bytes.NewBuffer(resp).String()
	responseWriter := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(stat)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, responseString)
	}
	return httptest.NewServer(http.HandlerFunc(responseWriter))
}

// responseIs returns true iff the response status code and body are identical
// to the provided expectations.
func responseIs(resp *http.Response, status int, expected []byte) bool {
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
