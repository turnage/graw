// Package reap provides an high level interface for Reddit HTTP requests.
package reap

import (
	"net/http"
	"net/url"

	"github.com/turnage/graw/internal/client"
	"github.com/turnage/graw/reddit"
)

var (
	// scheme is a map of TLS=[true|false] to the scheme for that setting.
	scheme = map[bool]string{
		true:  "https",
		false: "http",
	}
	formEncoding = map[string][]string{
		"content-type": {"application/x-www-form-urlencoded"},
	}
)

// Harvest contains Reddit elements yielded in a reaping.
type Harvest struct {
	Comments []*reddit.Comment
	Posts    []*reddit.Post
	Messages []*reddit.Message
}

type Config struct {
	Client   client.Client
	Parser   reddit.Parser
	Hostname string
	TLS      bool
}

// Reaper is a high level api for Reddit HTTP requests.
type Reaper interface {
	// Reap executes a GET request to Reddit and returns the elements from
	// the endpoint.
	Reap(path string, values map[string]string) (Harvest, error)
	// Sow executes a POST request to Reddit.
	Sow(path string, values map[string]string) error
}

type reaper struct {
	cli      client.Client
	parser   reddit.Parser
	hostname string
	scheme   string
}

// New returns a new Reaper.
func New(c Config) Reaper {
	return &reaper{
		cli:      c.Client,
		parser:   c.Parser,
		hostname: c.Hostname,
		scheme:   scheme[c.TLS],
	}
}

func (r *reaper) Reap(path string, values map[string]string) (Harvest, error) {
	harvest := Harvest{}
	resp, err := r.cli.Do(
		&http.Request{
			Method: "GET",
			URL:    r.url(path, values),
			Host:   r.hostname,
		},
	)
	if err != nil {
		return harvest, err
	}

	comments, posts, messages, err := r.parser.Parse(resp)
	harvest = Harvest{
		Comments: comments,
		Posts:    posts,
		Messages: messages,
	}

	return harvest, err
}

func (r *reaper) Sow(path string, values map[string]string) error {
	_, err := r.cli.Do(
		&http.Request{
			Method: "POST",
			Header: formEncoding,
			Host:   r.hostname,
			URL:    r.url(path, values),
		},
	)

	return err
}

func (r *reaper) url(path string, values map[string]string) *url.URL {
	return &url.URL{
		Scheme:   r.scheme,
		Host:     r.hostname,
		Path:     path,
		RawQuery: r.formatValues(values).Encode(),
	}
}

func (r *reaper) formatValues(values map[string]string) url.Values {
	formattedValues := url.Values{}

	for key, value := range values {
		formattedValues[key] = []string{value}
	}

	return formattedValues
}
