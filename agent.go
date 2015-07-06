// Package graw is the Golang Reddit API Wrapper.
package graw

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/paytonturnage/graw/nface"
	"golang.org/x/oauth2"
)

const (
	// authURL is the url for authorization requests.
	authURL = "https://www.reddit.com/api/v1/access_token"
	// baseURL is the base url for all api calls.
	baseURL = "https://oauth.reddit.com/api"
)

// api urls must be defined such that baseURL + apiURL makes the full api call.
const (
	meURL = "/v1/me"
)

type Agent struct {
	// client manages the Agent's oauth Authorization.
	client *http.Client
	// userAgent is written the http header of all requests.
	userAgent string
}

// NewAgent returns an agent ready to use with reddit. An error is returned
// if the agent cannot authenticate.
//
// See https://github.com/reddit/reddit/wiki/OAuth2 for more information.
func NewAgent(userAgent, id, secret, user, pass string) (*Agent, error) {
	conf := &oauth2.Config{
		ClientID: id,
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			TokenURL: authURL,
		},
	}

	token, err := conf.PasswordCredentialsToken(oauth2.NoContext, user, pass)
	if err != nil {
		return nil, err
	}

	if !token.Valid() {
		return nil, errors.New("received invalid token")
	}

	return &Agent{
		client: conf.Client(oauth2.NoContext, token),
		userAgent: userAgent,
	}, nil
}

func (a *Agent) Me() (*Redditor, error) {
	resp := &redditorResponse{}
	err := nface.Exec(
		a.client,
		a.userAgent,
		&nface.Request{
			Action: nface.GET,
			BaseUrl: baseURL + meURL,
		},
		resp)
	return &resp.Redditor, err
}
