package client

import (
	"net/http"
)

// Client defines behavior for making http.Requests.
type Client interface {
	// Do executes an http.Request.
	Do(r *http.Request) (*http.Response, error)
}
