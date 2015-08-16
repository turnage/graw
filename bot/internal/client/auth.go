package client

import (
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	// defaultAccessToken is the initial access token for oauth2 tokens.
	defaultAccessToken = "access_token"
	// tokenURL is the url of reddit's oauth2 authorization service.
	tokenURL = "https://reddit.com/api/v1/access_token"
)

// build returns an http clientt that has built in oauth2 handling.
func build(id, secret, refresh string) *http.Client {
	cfg := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint:     oauth2.Endpoint{TokenURL: tokenURL},
	}

	// For reddit, a permanent refresh token lasts forever. This initial
	// state tells oauth2 to use the refresh token to get a real access
	// token.
	token := &oauth2.Token{
		AccessToken:  defaultAccessToken,
		RefreshToken: refresh,
		Expiry:       time.Now().Add(-1 * time.Minute),
	}

	return cfg.Client(oauth2.NoContext, token)
}
