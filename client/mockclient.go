package client

import (
	"net/http"
)

type mockClient struct {
	resp *http.Response
	err  error
}

// Do returns the preconfigured response and error in the mock client.
func (m *mockClient) Do(r *http.Request) (*http.Response, error) {
	return m.resp, m.err
}
