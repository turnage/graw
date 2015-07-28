package testutil

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
)

type bytesCloser struct {
	buffer *bytes.Buffer
	err    error
}

func (b bytesCloser) Read(p []byte) (int, error) {
	if b.err != nil {
		return 0, b.err
	}

	return b.buffer.Read(p)
}

func (b bytesCloser) Close() error {
	return nil
}

// NewReadCloser returns a an io.ReadCloser with the content provided, which
// will return the err provided on calls to Read(). If err is nil, the error
// from bytes.Buffer will propogate up unadultered.
func NewReadCloser(content string, err error) io.ReadCloser {
	return &bytesCloser{
		buffer: bytes.NewBufferString(content),
		err:    err,
	}
}

// NewServerFromResponse returns an httptest.Server that always responds with
// response and status.
func NewServerFromResponse(stat int, resp []byte) *httptest.Server {
	responseString := bytes.NewBuffer(resp).String()
	responseWriter := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(stat)
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
