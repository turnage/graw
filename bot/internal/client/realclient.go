package client

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

const (
	// rateLimit is the wait time between queries.
	rateLimit = 3 * time.Second
)

type client struct {
	// agent is the client's User-Agent in http requests.
	agent string
	// id is the bot's OAuth2 client id.
	id string
	// secret is the bot's OAuth2 client secret.
	secret string
	// user is the bot's username on reddit.
	user string
	// pass is the bot's password on reddit.
	pass string

	// authMu protects authentication fields.
	authMu sync.Mutex
	// cli is the authenticated client to execute requests with.
	cli *http.Client
	// token is the OAuth2 token cli uses to authenticate.
	token *oauth2.Token

	// rateMu protects rate limiting fields.
	rateMu sync.Mutex
	// nextReq is the time at which it is ok to make the next request.
	nextReq time.Time
}

// Do wraps the execution of http requests. It updates authentications and rate
// limits requests to Reddit to comply with the API rules. It returns the
// response body.
func (c *client) Do(r *http.Request) (io.ReadCloser, error) {
	if time.Now().Before(c.nextReq) {
		<-time.After(c.nextReq.Sub(time.Now()))
	}
	c.nextReq = time.Now().Add(rateLimit)
	if !c.token.Valid() {
		var err error
		c.cli, c.token, err = build(c.id, c.secret, c.user, c.pass)
		if err != nil {
			return nil, err
		}
	}
	return c.exec(r)
}

// exec executes an http request and returns the response body.
func (c *client) exec(r *http.Request) (io.ReadCloser, error) {
	resp, err := c.doRaw(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response code: %d\n"+
			"request was: %v\n",
			resp.StatusCode,
			r)
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("no body in response")
	}

	return resp.Body, nil
}

// doRaw executes an http Request using an authenticated client, and the configured
// user agent.
func (c *client) doRaw(r *http.Request) (*http.Response, error) {
	if r.Header == nil {
		r.Header = make(http.Header)
	}
	r.Header.Add("User-Agent", c.agent)
	return c.cli.Do(r)
}
