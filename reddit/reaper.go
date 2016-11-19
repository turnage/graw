package reddit

import (
	"net/http"
	"net/url"
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

type reaperConfig struct {
	client   client
	parser   parser
	hostname string
	tls      bool
}

// reaper is a high level api for Reddit HTTP requests.
type reaper interface {
	// reap executes a GET request to Reddit and returns the elements from
	// the endpoint.
	reap(path string, values map[string]string) (Harvest, error)
	// sow executes a POST request to Reddit.
	sow(path string, values map[string]string) error
}

type reaperImpl struct {
	cli      client
	parser   parser
	hostname string
	scheme   string
}

func newReaper(c reaperConfig) reaper {
	return &reaperImpl{
		cli:      c.client,
		parser:   c.parser,
		hostname: c.hostname,
		scheme:   scheme[c.tls],
	}
}

func (r *reaperImpl) reap(path string, values map[string]string) (Harvest, error) {
	resp, err := r.cli.Do(
		&http.Request{
			Method: "GET",
			URL:    r.url(path, values),
			Host:   r.hostname,
		},
	)
	if err != nil {
		return Harvest{}, err
	}

	comments, posts, messages, err := r.parser.parse(resp)
	return Harvest{
		Comments: comments,
		Posts:    posts,
		Messages: messages,
	}, err
}

func (r *reaperImpl) sow(path string, values map[string]string) error {
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

func (r *reaperImpl) url(path string, values map[string]string) *url.URL {
	return &url.URL{
		Scheme:   r.scheme,
		Host:     r.hostname,
		Path:     path,
		RawQuery: r.formatValues(values).Encode(),
	}
}

func (r *reaperImpl) formatValues(values map[string]string) url.Values {
	formattedValues := url.Values{}

	for key, value := range values {
		formattedValues[key] = []string{value}
	}

	return formattedValues
}
