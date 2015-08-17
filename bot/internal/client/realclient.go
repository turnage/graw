package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type client struct {
	agent  string
	id     string
	secret string
	user   string
	pass   string
	token  *oauth2.Token
	cli    *http.Client
}

func (c *client) Do(r *http.Request, out interface{}) error {
	if !c.token.Valid() {
		var err error
		c.cli, c.token, err = build(c.id, c.secret, c.user, c.pass)
		if err != nil {
			return err
		}
	}
	return c.exec(r, out)
}

func (c *client) exec(r *http.Request, out interface{}) error {
	rawResp, err := c.doRaw(r)
	if err != nil {
		return err
	}

	if rawResp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response code: %d\n"+
			"request was: %v\n",
			rawResp.StatusCode,
			r)
	}

	defer func() {
		if rawResp.Body != nil {
			rawResp.Body.Close()
		}
	}()

	decoder := json.NewDecoder(rawResp.Body)
	return decoder.Decode(out)
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
