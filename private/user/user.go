package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/paytonturnage/graw/private/auth"
	"github.com/paytonturnage/graw/private/client"
	"github.com/paytonturnage/redditproto"
)

const (
	// authURL is the url for authorization requests.
	authURL = "https://www.reddit.com/api/v1/access_token"
)

// netUser implements User with functions that control the real network.
type User struct {
	// agent is the user controller (bot/script) user agent.
	agent string
	// authorizer handles authentication with reddit
	authorizer auth.Authorizer
	// client executes all network requests.
	client client.Client
}

// New returns an authenticated reddit user which can be controlled to make
// requests and interact with reddit.
func New(agent *redditproto.UserAgent) *User {
	return &User{
		agent: agent.GetUserAgent(),
		authorizer: auth.NewOAuth2Authorizer(
			agent.GetClientId(),
			agent.GetClientSecret(),
			agent.GetUsername(),
			agent.GetPassword(),
		)}
}

// Auth identifies as the user to the Reddit servers.
func (u *User) Auth() error {
	var err error
	u.client, err = u.authorizer.Client(authURL)
	return err
}

func (u *User) Exec(req *http.Request, resp interface{}) error {
	rawResp, err := u.ExecRaw(req)
	if err != nil {
		return err
	}

	if rawResp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code in response")
	}

	if rawResp.Body == nil {
		return fmt.Errorf("no body in response")
	}
	defer rawResp.Body.Close()

	buffer, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(buffer, resp)
}

func (u *User) ExecRaw(r *http.Request) (*http.Response, error) {
	return u.client.Do(r)
}
