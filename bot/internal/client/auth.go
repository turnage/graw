package client

import (
	"net/http"

	"golang.org/x/oauth2"
)

const (
	// tokenURL is the url of reddit's oauth2 authorization service.
	tokenURL = "https://www.reddit.com/api/v1/access_token"
)

// build returns an http clientt that has built in oauth2 handling.
func build(id, secret, user, pass string) (*http.Client, *oauth2.Token, error) {
	cfg := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint:     oauth2.Endpoint{TokenURL: tokenURL},
	}
	token, err := cfg.PasswordCredentialsToken(oauth2.NoContext, user, pass)
	return cfg.Client(oauth2.NoContext, token), token, err
}
