// Package client reaps http Requests which require OAuth2 authorization.
package client

import (
	"bytes"
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var (
	PermissionDeniedErr = fmt.Errorf("unauthorized access to endpoint")
	BusyErr             = fmt.Errorf("service is busy right now")
	RateLimitErr        = fmt.Errorf("service rate limiting requests")
)

// Config holds all the information needed to define Client behavior, such as
// who the client will identify as externally and where to authorize.
type Config struct {
	// Agent is the user agent set in all requests made by the Client.
	Agent string

	// ID and Secret are used to claim an OAuth2 grant the users are
	// previously authorized.
	ID     string
	Secret string

	// Username and Password are used to authorize with the endpoint.
	Username string
	Password string

	// TokenURL is the url of the token request location for OAuth2.
	TokenURL string
}

// Client executes http Requests and invisibly handles OAuth2 authorization.
type Client interface {
	Do(*http.Request) ([]byte, error)
}

type client struct {
	cli *http.Client
}

func (r *client) Do(req *http.Request) ([]byte, error) {
	resp, err := r.cli.Do(req)
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
	default:
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// New returns a new Client using the given user to make requests.
func New(c Config) (Client, error) {
	transport := &agentForwarder{agent: c.Agent}
	cli := &http.Client{Transport: transport}
	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, cli)
	cfg := &oauth2.Config{
		ClientID:     c.ID,
		ClientSecret: c.Secret,
		Endpoint:     oauth2.Endpoint{TokenURL: c.TokenURL},
		Scopes: []string{
			"identity",
			"read",
			"privatemessages",
			"submit",
			"history",
		},
	}

	token, err := cfg.PasswordCredentialsToken(ctx, c.Username, c.Password)
	return &client{cfg.Client(ctx, token)}, err
}
