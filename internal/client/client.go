// Package client reaps http Requests which require OAuth2 authorization.
package client

import (
	"bytes"
	"fmt"
	"net/http"
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

	// If all fields in App are set, this client will attempt to identify as
	// a registered Reddit app using the credentials.
	App App
}

// Client executes http Requests and invisibly handles OAuth2 authorization.
type Client interface {
	Do(*http.Request) ([]byte, error)
}

type base struct {
	cli *http.Client
}

func (b *base) Do(req *http.Request) ([]byte, error) {
	resp, err := b.cli.Do(req)
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
	if c.App.configured() {
		return newAppClient(c)
	}

	return &base{clientWithAgent(c.Agent)}, nil
}
