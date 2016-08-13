package client

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var (
	// TokenURL is the url of reddit's oauth2 authorization service.
	TokenURL = "https://www.reddit.com/api/v1/access_token"
	// TestMode is a flag that indicates whether graw is in test mode.
	TestMode = false
)

type oauthTransport struct {
	http.Transport
	agent string
}

func (o *oauthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("User-Agent", o.agent)
	return o.Transport.RoundTrip(r)
}

// build returns an http client that has built in oauth2 handling.
func build(agent, id, secret, user, pass string) (*http.Client, *oauth2.Token, error) {
	cfg := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint:     oauth2.Endpoint{TokenURL: TokenURL},
		Scopes: []string{
			"identity",
			"read",
			"privatemessages",
			"submit",
			"history",
		},
	}

	if TestMode {
		return buildTestClient(cfg, user, pass)
	}

	return buildProductionClient(cfg, agent, user, pass)
}

// buildProductionClient returns a client equipped to make requests to
// a production Reddit instance.
func buildProductionClient(
	cfg *oauth2.Config,
	agent,
	user,
	pass string,
) (
	*http.Client,
	*oauth2.Token,
	error,
) {
	cli := http.DefaultClient
	cli.Transport = &oauthTransport{agent: agent}
	token, err := cfg.PasswordCredentialsToken(
		context.WithValue(oauth2.NoContext, oauth2.HTTPClient, cli),
		user,
		pass,
	)
	return cfg.Client(oauth2.NoContext, token), token, err
}

// buildTestClient returns a client equipped to make requests to a production
// Reddit instance.
func buildTestClient(
	cfg *oauth2.Config,
	user,
	pass string,
) (
	*http.Client,
	*oauth2.Token,
	error,
) {
	naiveTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	naiveContext := context.WithValue(
		oauth2.NoContext,
		oauth2.HTTPClient,
		&http.Client{
			Transport: naiveTransport,
		},
	)

	token, err := cfg.PasswordCredentialsToken(naiveContext, user, pass)
	client := cfg.Client(oauth2.NoContext, token)
	client.Transport.(*oauth2.Transport).Base = naiveTransport
	return client, token, err
}
