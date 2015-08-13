package graw

import (
	"net/http"
)

// netClient implements Client and provides http.Client and OAuth functionality.
type netClient struct {
	// client is the network client.
	Client *http.Client
}

func (n *netClient) Do(r *http.Request) (*http.Response, error) {
	return n.Client.Do(r)
}
