package client

import (
	"net/http"
)

type client struct {
	useragent string
	cli       *http.Client
}

// Do executes an http Request using a configured client, and ensuring that the
// configured user agent is set.
func (c *client) Do(r *http.Request) (*http.Response, error) {
	if r.Header == nil {
		r.Header = make(http.Header)
	}
	r.Header.Add("User-Agent", c.useragent)
	return c.cli.Do(r)
}
