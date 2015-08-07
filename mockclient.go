package graw

import (
	"net/http"
)

type mockClient struct {
	// response will be provided by all calls to Do.
	response *http.Response
	// err will be provided by all calls to Do.
	err error
}

// do mocks execution of an http.Request, instead returning the mockClient
// instance's preset response and error values.
func (m *mockClient) do(r *http.Request) (*http.Response, error) {
	return m.response, m.err
}
