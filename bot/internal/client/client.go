// Package client manages all requests made to reddit.
package client

import (
	"fmt"
	"net/http"
)

// Client implementations make requests.
type Client interface {
	// Do executes a request and manages authenticating with and identifying
	// to the server.
	Do(r *http.Request) (*http.Response, error)
}

// New returns a new Client from a user agent file. This user agent file is
// expected to have been generated with "graw grant".
func New(filename string) (Client, error) {
	agent, err := load(filename)
	if err != nil {
		return nil, err
	}

	if agent.GetRefreshToken() == "" {
		return nil, fmt.Errorf("no refresh token; see graw grant")
	}

	return &client{cli: build(
		agent.GetUserAgent(),
		agent.GetClientId(),
		agent.GetRefreshToken(),
	)}, nil
}

// NewMock returns a mock client that returns canned values for calls to Do().
func NewMock(r *http.Response, err error) Client {
	return &mockClient{r, err}
}
