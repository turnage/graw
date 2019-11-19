package reddit

import (
	"bytes"
	"errors"
	"fmt"
	"log"
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
	autoRefresh *autoRefresh
	cfg         clientConfig
	cli         *http.Client
	app         *App
}

func (a *appClient) AutoRefresh() error {
	if a.autoRefresh == nil {
		return errors.New("autoRefresh not set")
	}
	go a.autoRefresh.autoRefresh(a)
	return nil
}

func (a *appClient) Do(req *http.Request) ([]byte, error) {
	resp, err := a.cli.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusForbidden:
		return nil, PermissionDeniedErr
	case http.StatusServiceUnavailable:
		return nil, BusyErr
	case http.StatusTooManyRequests:
		return nil, RateLimitErr
	case http.StatusBadGateway:
		return nil, GatewayErr
	case http.StatusGatewayTimeout:
		return nil, GatewayTimeoutErr
	default:
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (a *appClient) authorize() error {
	a.cli = clientWithAgent(a.cfg.agent)
	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, a.cli)

	if a.cfg.app.Username == "" || a.cfg.app.Password == "" {
		a.cli = a.clientCredentialsClient(ctx)
		return nil
	}

	cfg := &oauth2.Config{
		ClientID:     a.cfg.app.ID,
		ClientSecret: a.cfg.app.Secret,
		Endpoint:     oauth2.Endpoint{TokenURL: a.cfg.app.tokenURL},
		Scopes:       oauthScopes,
	}

	token, err := cfg.PasswordCredentialsToken(
		ctx,
		a.cfg.app.Username,
		a.cfg.app.Password,
	)

	if err != nil {
		return err
	}

	log.Printf("access_token generated: %s\n", token.AccessToken)
	a.cli = cfg.Client(ctx, token)
	a.autoRefresh.setRefreshTimerFromToken(token)
	return nil
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
		cli:         clientWithAgent(c.agent),
		cfg:         c,
		autoRefresh: newAutoRefresh(),
	}
	return a, a.authorize()
}
