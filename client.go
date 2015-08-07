package graw

import (
	"net/http"
)

// client defines behavior for making http.Requests.
type client interface {
	// do executes an http.Request.
	do(r *http.Request) (*http.Response, error)
}
