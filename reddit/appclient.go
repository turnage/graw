package reddit

import (
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type appClient struct {
	baseClient
	cfg   clientConfig
	cli   *http.Client
	token *oauth2.Token
}

func (a *appClient) Do(req *http.Request) ([]byte, error) {
	if a.token == nil || !a.token.Valid() {
		if err := a.authorize(); err != nil {
			return nil, err
		}
	}

	return a.baseClient.Do(req)
}

func (a *appClient) authorize() error {
	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, a.cli)
	cfg := &oauth2.Config{
		ClientID:     a.cfg.app.ID,
		ClientSecret: a.cfg.app.Secret,
		Endpoint:     oauth2.Endpoint{TokenURL: a.cfg.app.tokenURL},
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
		a.cfg.app.Username,
		a.cfg.app.Password,
	)

	a.token = token
	a.baseClient.cli = cfg.Client(ctx, token)
	return err
}

func newAppClient(c clientConfig) (*appClient, error) {
	a := &appClient{
		cli: clientWithAgent(c.agent),
		cfg: c,
	}
	return a, a.authorize()
}
