package client

import (
	"net/http"
)

// netClient implements Client and provides http.Client and OAuth functionality.
type netClient struct {
	// client is the network client.
	client *http.Client
}

// NewNetClient provides a Client that executes http.Requests over the real
// network.
func NewNetClient(c *http.Client) Client {
	return &netClient{client: c}
}

func (n *netClient) Do(r *http.Request) (*http.Response, error) {
	return n.client.Do(r)
}
