package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type client struct {
	useragent string
	cli       *http.Client
}

func (c *client) Do(r *http.Request, out interface{}) error {
	rawResp, err := c.doRaw(r)
	if err != nil {
		return err
	}

	if rawResp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response code: %d", rawResp.StatusCode)
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
	r.Header.Add("User-Agent", c.useragent)
	return c.cli.Do(r)
}
