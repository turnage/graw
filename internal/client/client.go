package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

const (
	// rateLimit is the wait time between requests to Reddit.
	rateLimit = 2 * time.Second
)

type Client struct {
	// agent is the Client's User-Agent in http requests.
	agent string
	// id is the bot's OAuth2 Client id.
	id string
	// secret is the bot's OAuth2 Client secret.
	secret string
	// user is the bot's username on reddit.
	user string
	// pass is the bot's password on reddit.
	pass string

	// authMu protects authentication fields.
	authMu sync.Mutex
	// cli is the authenticated Client to execute requests with.
	cli *http.Client
	// token is the OAuth2 token cli uses to authenticate.
	token *oauth2.Token

	// rateMu protects rate limiting fields.
	rateMu sync.Mutex
	// nextReq is the time at which it is ok to make the next request.
	nextReq time.Time
}

// New returns a new Client from a user agent file.
func New(filename string) (*Client, error) {
	agent, err := load(filename)
	if err != nil {
		return nil, err
	}

	return &Client{
		agent:   agent.GetUserAgent(),
		id:      agent.GetClientId(),
		secret:  agent.GetClientSecret(),
		user:    agent.GetUsername(),
		pass:    agent.GetPassword(),
		nextReq: time.Now(),
	}, nil
}

// Do wraps the execution of http requests. It updates authentications and rate
// limits requests to Reddit to comply with the API rules. It returns the
// response body.
func (c *Client) Do(r *http.Request) ([]byte, error) {
	c.rateRequest()

	if !c.token.Valid() {
		var err error
		c.cli, c.token, err = build(c.agent, c.id, c.secret, c.user, c.pass)
		if err != nil {
			return nil, err
		}
	}

	body, err := c.exec(r)
	if err != nil {
		return nil, err
	}

	return responseBytes(body)
}

// exec executes an http request and returns the response body.
func (c *Client) exec(r *http.Request) (io.ReadCloser, error) {
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

// doRaw executes an http Request using an authenticated Client, and the configured
// user agent.
func (c *Client) doRaw(r *http.Request) (*http.Response, error) {
	if r.Header == nil {
		r.Header = make(http.Header)
	}
	r.Header.Add("User-Agent", c.agent)
	return c.cli.Do(r)
}

// rateRequest blocks until the rate limits have been abided by.
func (c *Client) rateRequest() {
	c.rateMu.Lock()
	defer c.rateMu.Unlock()

	if time.Now().After(c.nextReq) {
		c.nextReq = time.Now().Add(rateLimit)
		return
	}

	currentReq := c.nextReq
	c.nextReq = currentReq.Add(rateLimit)
	<-time.After(currentReq.Sub(time.Now()))
}

// responseBytes returns a slice of bytes from a response body.
func responseBytes(response io.ReadCloser) ([]byte, error) {
	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(response); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
