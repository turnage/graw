package graw

import (
	"net/http"
)

// client defines behavior for making http.Requests.
type client interface {
	// Do executes an http.Request.
	Do(r *http.Request) (*http.Response, error)
}
