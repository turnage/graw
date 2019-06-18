package reddit

import (
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var oauthScopes = []string{
	"identity",
	"read",
	"privatemessages",
	"submit",
	"history",
}

type appClient struct {
	baseClient
	cfg clientConfig
	cli *http.Client
}

func (a *appClient) Do(req *http.Request) ([]byte, error) {
	return a.baseClient.Do(req)
}

func (a *appClient) authorize() error {
	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, a.cli)

	if a.cfg.app.unauthenticated() {
		a.baseClient.cli = a.clientCredentialsClient(ctx)
		return nil
	}

	cfg := &oauth2.Config{
		ClientID:     a.cfg.app.ID,
		ClientSecret: a.cfg.app.Secret,
		Endpoint:     oauth2.Endpoint{TokenURL: a.cfg.app.tokenURL},
		Scopes:       oauthScopes,
	}

	var token *oauth2.Token
	var err error

	if a.cfg.app.Token != nil {
		token = a.cfg.app.Token
		err = nil
	} else {
		token, err = cfg.PasswordCredentialsToken(
			ctx,
			a.cfg.app.Username,
			a.cfg.app.Password,
		)
	}

	a.baseClient.cli = cfg.Client(ctx, token)

	return err
}

func (a *appClient) clientCredentialsClient(ctx context.Context) *http.Client {
	cfg := &clientcredentials.Config{
		ClientID:     a.cfg.app.ID,
		ClientSecret: a.cfg.app.Secret,
		TokenURL:     a.cfg.app.tokenURL,
		Scopes:       oauthScopes,
	}

	return cfg.Client(ctx)
}

func newAppClient(c clientConfig) (*appClient, error) {
	a := &appClient{
		cli: clientWithAgent(c.agent),
		cfg: c,
	}
	return a, a.authorize()
}
