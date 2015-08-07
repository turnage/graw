package graw

import (
	"golang.org/x/oauth2"
	"net/http"
)

func oauth(id, secret, user, pass, url string) (*http.Client, error) {
	cfg := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			TokenURL: url,
		},
	}

	token, err := cfg.PasswordCredentialsToken(oauth2.NoContext, user, pass)
	return cfg.Client(oauth2.NoContext, token), err
}
