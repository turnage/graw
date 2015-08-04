package graw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/paytonturnage/graw/internal/auth"
	"github.com/paytonturnage/graw/internal/client"
	"github.com/paytonturnage/graw/internal/ratelimiter"
	"github.com/paytonturnage/graw/internal/request"
	"github.com/paytonturnage/redditproto"
)

const (
	// authURL is the url for authorization requests.
	authURL = "https://www.reddit.com/api/v1/access_token"
	// maxQueriesPerMinute
	maxQueriesPerMinute = 60
)

// User is a reddit user.
type User struct {
	// agent is the user controller (bot/script) user agent.
	agent string
	// authorizer handles authentication with reddit
	authorizer auth.Authorizer
	// client executes all network requests.
	client client.Client
	// limiter limits the queries made per api window.
	limiter ratelimiter.RateLimiter
}

// NewUser returns an authenticated reddit user which can be controlled to make
// requests and interact with reddit.
func NewUser(agent *redditproto.UserAgent) *User {
	return &User{
		agent: agent.GetUserAgent(),
		authorizer: auth.NewOAuth2Authorizer(
			agent.GetClientId(),
			agent.GetClientSecret(),
			agent.GetUsername(),
			agent.GetPassword(),
		),
		limiter: ratelimiter.NewTimeRateLimiter(
			time.Minute,
			maxQueriesPerMinute,
		),
	}
}

// NewUserFromFile returns an authenticated reddit user which can be controlled
// to make requests and interact with reddit from a user agent protobuf file.
func NewUserFromFile(filename string) (*User, error) {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	agent := &redditproto.UserAgent{}
	if err := proto.UnmarshalText(
		bytes.NewBuffer(buffer).String(),
		agent,
	); err != nil {
		return nil, err
	}
	return NewUser(agent), err
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
	if u.limiter != nil {
		u.limiter.Request(true)
	}
	return u.client.Do(r)
}

// Me returns an account data structure that represents the logged in user.
func (u *User) Me() (*redditproto.Account, error) {
	account := &redditproto.Account{}
	req, err := request.New(
		"GET",
		"https://oauth.reddit.com/api/v1/me",
		nil,
	)
	if err != nil {
		return nil, err
	}

	return account, u.Exec(req, account)
}

func (u *User) Scrape(sub, sort, after, before string,
	lim int) ([]*redditproto.Link, error) {
	response := &struct {
		Data struct {
			Children []struct {
				Data *redditproto.Link
			}
		}
	}{}
	req, err := request.New(
		"GET",
		fmt.Sprintf("https://oauth.reddit.com/r/%s/%s.json", sub, sort),
		&url.Values{
			"limit":  []string{strconv.Itoa(lim)},
			"after":  []string{after},
			"before": []string{before},
		},
	)
	if err != nil {
		return nil, err
	}

	err = u.Exec(req, response)
	if err != nil {
		return nil, err
	}

	links := make([]*redditproto.Link, len(response.Data.Children))
	for i, child := range response.Data.Children {
		links[i] = child.Data
	}

	return links, nil
}
