package graw

import (
	"net/http"
)

// netClient implements Client and provides http.Client and OAuth functionality.
type netClient struct {
	// client is the network client.
	client *http.Client
}

func (n *netClient) do(r *http.Request) (*http.Response, error) {
	return n.client.Do(r)
}
