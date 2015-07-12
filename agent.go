// Package graw is the Golang Reddit API Wrapper.
package graw

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/paytonturnage/graw/api"
	"github.com/paytonturnage/graw/nface"
	"github.com/paytonturnage/graw/data"
	"golang.org/x/oauth2"
)

const (
	// authURL is the url for authorization requests.
	authURL = "https://www.reddit.com/api/v1/access_token"
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

// NewAgentFromFile calls NewAgent with auth information read from a
// protobuffer file. See useragent.protobuf.example.
func NewAgentFromFile(filename string) (*Agent, error) {
	agentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	agentText := bytes.NewBuffer(agentBytes)
	agent := &data.UserAgent{}
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
func (a *Agent) Me() (*data.Redditor, error) {
	resp := &data.Redditor{}
	err := a.client.Do(api.MeRequest(), resp)
	return resp, err
}

// MeKarma wraps /v1/me/karma. See
// https://www.reddit.com/dev/api#GET_api_v1_me_karma 
func (a *Agent) MeKarma() (*data.KarmaList, error) {
	resp := &data.KarmaList{}
	err := a.client.Do(api.MeKarmaRequest(), resp)
	return resp, err
}
