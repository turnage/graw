package graw

import (
	"net/http"
)

type mockClient struct {
	// response will be provided by all calls to Do.
	Response *http.Response
	// err will be provided by all calls to Do.
	Err error
}

// do mocks execution of an http.Request, instead returning the mockClient
// instance's preset response and error values.
func (m *mockClient) Do(r *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}
