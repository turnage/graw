// Package reaper reaps http Requests which require OAuth2 authorization.
package reaper

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

// Config holds all the information needed to define Reaper behavior, such as
// who the reaper will identify as externally and where to authorize.
type Config struct {
	// Agent is the user agent set in all requests made by the Reaper.
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

// Reaper executes http Requests and invisibly handles OAuth2 authorization.
type Reaper interface {
}

type reaper struct {
	cli *http.Client
}

func (r *reaper) reap(req *http.Request) (string, error) {
	resp, err := r.cli.Do(req)
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusForbidden:
		return "", PermissionDeniedErr
	case http.StatusServiceUnavailable:
		return "", BusyErr
	case http.StatusTooManyRequests:
		return "", RateLimitErr
	default:
		return "", fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// New returns a new Reaper using the given user to reap requests.
func New(c Config) (Reaper, error) {
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
	return &reaper{cfg.Client(ctx, token)}, err
}
