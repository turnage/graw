package client

import (
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type appClient struct {
	base
	cli   *http.Client
	cfg   Config
	token *oauth2.Token
}

func (a *appClient) Do(req *http.Request) ([]byte, error) {
	if a.token == nil || !a.token.Valid() {
		if err := a.authorize(); err != nil {
			return nil, err
		}
	}

	return a.base.Do(req)
}

func (a *appClient) authorize() error {
	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, a.cli)
	cfg := &oauth2.Config{
		ClientID:     a.cfg.App.ID,
		ClientSecret: a.cfg.App.Secret,
		Endpoint:     oauth2.Endpoint{TokenURL: a.cfg.App.TokenURL},
		Scopes: []string{
			"identity",
			"read",
			"privatemessages",
			"submit",
			"history",
		},
	}

	token, err := cfg.PasswordCredentialsToken(
		ctx,
		a.cfg.App.Username,
		a.cfg.App.Password,
	)

	a.token = token
	a.base.cli = cfg.Client(ctx, token)
	return err
}

func newAppClient(c Config) (*appClient, error) {
	a := &appClient{
		cli: clientWithAgent(c.Agent),
		cfg: c,
	}
	return a, a.authorize()
}
