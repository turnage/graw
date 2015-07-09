// Package graw is the Golang Reddit API Wrapper.
package graw

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
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

// Agent wraps the reddit api; all api calls go through Agent.
type Agent struct {
	// client manages the Agent's connection to reddit.
	client *nface.Client
}

// NewAgent returns an agent ready to use with reddit. An error is returned
// if the agent cannot authenticate.
//
// See https://github.com/reddit/reddit/wiki/OAuth2 for more information.
func NewAgent(userAgent, id, secret, user, pass string) (*Agent, error) {
	conf := &oauth2.Config{
		ClientID:     id,
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

	return &Agent{client: nface.NewClient(
		conf.Client(oauth2.NoContext, token), userAgent)}, nil
}

// NewAgentFromFile returns an agent with auth information read from a
// protobuffer file. See the UserAgent message type in graw.proto for what
// fields to provide in the file.
func NewAgentFromFile(filename string) (*Agent, error) {
	agentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	agentText := bytes.NewBuffer(agentBytes)
	agent := &UserAgent{}
	if err := proto.UnmarshalText(agentText.String(), agent); err != nil {
		return nil, err
	}

	return NewAgent(
		agent.GetUserAgent(),
		agent.GetClientId(),
		agent.GetClientSecret(),
		agent.GetUsername(),
		agent.GetPassword())
}

// Me wraps /v1/me. See https://www.reddit.com/dev/api#GET_api_v1_me
func (a *Agent) Me() (*Redditor, error) {
	resp := &Redditor{}
	err := a.client.Do(&nface.Request{
		Action:  nface.GET,
		BaseUrl: baseURL + meURL,
	}, resp)
	return resp, err
}
