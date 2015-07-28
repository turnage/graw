package client

import (
	"net/http"
)

type mockClient struct {
	// response will be provided by all calls to Do.
	response *http.Response
	// err will be provided by all calls to Do.
	err error
}

// NewMockClient provides an implementation of Client that does not use the real
// network at all; all calls to Do will provide the response and error value
// passed to this function.
func NewMockClient(r *http.Response, err error) Client {
	return &mockClient{response: r, err: err}
}

// Do mocks execution of an http.Request, instead returning the mockClient
// instance's preset response and error values.
func (m *mockClient) Do(r *http.Request) (*http.Response, error) {
	return m.response, m.err
}
