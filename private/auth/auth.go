// Package auth handles authentication with Reddit.
package auth

import (
	"net/http"
)

// Authorizer defines the behaviors of an authorizer method (e.g. OAuth2).
type Authorizer interface {
	// Client returns an http.Client that automatically handles
	// authorization on all requests made using it. authURL is the url of
	// the authorization server.
	Client(authURL string) (*http.Client, error)
}
