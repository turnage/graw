// Package client manages all requests made to reddit.
package client

import (
	"net/http"
)

// Client implementations make requests.
type Client interface {
	// Do executes a request to reddit and interprets the response into out.
	Do(r *http.Request, out interface{}) error
}

// New returns a new Client from a user agent file.
func New(filename string) (Client, error) {
	agent, err := load(filename)
	if err != nil {
		return nil, err
	}

	return &client{
		agent:  agent.GetUserAgent(),
		id:     agent.GetClientId(),
		secret: agent.GetClientSecret(),
		user:   agent.GetUsername(),
		pass:   agent.GetPassword(),
	}, nil
}

// NewMock returns a mock client that will parse the canned response string into
// the output interface for all calls to Do().
func NewMock(response string) Client {
	return &mockClient{response: []byte(response)}
}
