// Package graw is the Golang Reddit API Wrapper.
package graw

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/paytonturnage/graw/data"
	"github.com/paytonturnage/graw/nface"
	"golang.org/x/oauth2"
)

const (
	// authURL is the url for authorization requests.
	authURL = "https://www.reddit.com/api/v1/access_token"
	// baseURL is the base url for all api calls.
	baseURL = "https://oauth.reddit.com/api"
)

// Reddit api call urls.
const (
	// meURL is the url exension for the /v1/me api call.
	meURL = "/v1/me"
	// meKarmaURL is the url extension /v1/me/karma api call.
	meKarmaURL = "/v1/me/karma"
)

// Graw wraps the reddit api; all api calls go through Graw.
type Graw struct {
	// client manages the Graw's connection to reddit.
	client *nface.Client
}

// NewGraw returns an graw ready to use with reddit. An error is returned
// if the graw cannot authenticate.
//
// See https://github.com/reddit/reddit/wiki/OAuth2 for more information.
func NewGraw(userAgent, id, secret, user, pass string) (*Graw, error) {
	return newGrawFromUserAgent(newUserAgent(userAgent, id, secret, user, pass))
}

// NewGrawFromFile calls NewAgent with auth information read from a
// protobuffer file. See usergraw.protobuf.example.
func NewGrawFromFile(filename string) (*Graw, error) {
	userAgent, err := newUserAgentFromFile(filename)
	if err != nil {
		return nil, err
	}

	return newGrawFromUserAgent(userAgent)
}

// Me wraps /v1/me. See https://www.reddit.com/dev/api#GET_api_v1_me
func (a *Graw) Me() (*data.Account, error) {
	resp := &data.Account{}
	err := a.client.Do(&nface.Request{
		Action: nface.GET,
		URL:    meURL,
	}, resp)
	return resp, err
}

// MeKarma wraps /v1/me/karma. See
// https://www.reddit.com/dev/api#GET_api_v1_me_karma
func (a *Graw) MeKarma() ([]*data.SubredditKarma, error) {
	resp := &data.KarmaList{}
	err := a.client.Do(&nface.Request{
		Action: nface.GET,
		URL:    meKarmaURL,
	}, resp)
	return resp.GetData(), err
}

// newGrawFromUserAgent returns a new Graw derived from the identity in
// UserAgent.
func newGrawFromUserAgent(userAgent *data.UserAgent) (*Graw, error) {
	cli, err := newOAuthClient(userAgent, authURL)
	if err != nil {
		return nil, err
	}

	return &Graw{
		client: nface.NewClient(cli, userAgent.GetUserAgent(), baseURL),
	}, nil
}

// newUserAgent returns a new data.UserAgent containing the provided fields.
func newUserAgent(userAgent, id, secret, user, pass string) *data.UserAgent {
	return &data.UserAgent{
		UserAgent:    &userAgent,
		ClientId:     &id,
		ClientSecret: &secret,
		Username:     &user,
		Password:     &pass,
	}
}

// newUserAgent returns a new data.UserAgent from a protobuffer file.
func newUserAgentFromFile(filename string) (*data.UserAgent, error) {
	userAgentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	grawText := bytes.NewBuffer(userAgentBytes)
	userAgent := &data.UserAgent{}
	if err := proto.UnmarshalText(grawText.String(), userAgent); err != nil {
		return nil, err
	}

	return userAgent, nil
}

// newOAuthClient creates an OAuth http.Client that manages OAuth with reddit.
func newOAuthClient(user *data.UserAgent, auth string) (*http.Client, error) {
	conf := &oauth2.Config{
		ClientID:     user.GetClientId(),
		ClientSecret: user.GetClientSecret(),
		Endpoint: oauth2.Endpoint{
			TokenURL: auth,
		},
	}

	token, err := conf.PasswordCredentialsToken(
		oauth2.NoContext,
		user.GetUsername(),
		user.GetPassword())
	if err != nil {
		return nil, err
	}

	return conf.Client(oauth2.NoContext, token), nil
}
