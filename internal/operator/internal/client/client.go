// Package client manages all requests made to reddit.
package client

import (
	"io"
	"net/http"
	"time"
)

// Client implementations make requests.
type Client interface {
	// Do executes a request to reddit and returns the response body.
	Do(r *http.Request) (io.ReadCloser, error)
}

// New returns a new Client from a user agent file.
func New(filename string) (Client, error) {
	agent, err := load(filename)
	if err != nil {
		return nil, err
	}

	return &client{
		agent:   agent.GetUserAgent(),
		id:      agent.GetClientId(),
		secret:  agent.GetClientSecret(),
		user:    agent.GetUsername(),
		pass:    agent.GetPassword(),
		nextReq: time.Now(),
	}, nil
}

// NewMock returns a mock client that will return the canned values.
func NewMock(response string, err error) Client {
	return &mockClient{
		response: []byte(response),
		err:      err,
	}
}
