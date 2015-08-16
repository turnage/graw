// Package client manages all requests made to reddit.
package client

import (
	"net/http"
)

// Client implementations make requests.
type Client interface {
	// Do executes a request and manages authenticating with and identifying
	// to the server.
	Do(r *http.Request) (*http.Response, error)
}
