// Package operator makes api calls to Reddit.
package operator

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/turnage/graw/internal/operator/internal/client"
	"github.com/turnage/redditproto"
)

type Kind int

const (
	Comment = iota
	Link    = iota
)

const (
	// MaxLinks is the amount of posts reddit will return for a scrape
	// query.
	MaxLinks = 100
	// baseURL is the url all requests extend from.
	baseURL = "https://oauth.reddit.com"
	// oauth2Host is the hostname of Reddit's OAuth2 server.
	oauth2Host = "oauth.reddit.com"
)

var (
	formEncoding = map[string][]string{
		"content-type": []string{"application/x-www-form-urlencoded"},
	}
)

// Thing describes methods of Reddit's Thing class.
type Thing interface {
	// GetName returns the fullname of a Reddit Thing.
	GetName() string
}

// Operator makes api calls to Reddit.
type Operator interface {
	// Scrape fetches new reddit posts (see definition).
	Scrape(path, after, before string, limit uint, kind Kind) ([]Thing, error)
	// GetThing fetches a particular thing from reddit. Thing returns nil if
	// there is no such Thing.
	GetThing(id string, kind Kind) (Thing, error)
	// Thread fetches a post and its comment tree (see definition).
	Thread(permalink string) (*redditproto.Link, error)
	// Inbox fetches unread messages from the reddit inbox (see definition).
	Inbox() ([]*redditproto.Message, error)
	// MarkAsRead marks inbox items read (see definition).
	MarkAsRead(fullnames ...string) error
	// Reply replies to reddit item (see definition).
	Reply(parent, content string) error
	// Compose sends a private message to a user (see definition).
	Compose(user, subject, content string) error
	// Submit posts to Reddit (see definition).
	Submit(subreddit, kind, title, content string) error
}

// operator implements Operator.
type operator struct {
	cli client.Client
}

// New returns a new operator which uses agent as its identity. agent should be
// a filename of a file containing a UserAgent protobuffer.
func New(agent string) (Operator, error) {
	cli, err := client.New(agent)
	if err != nil {
		return nil, err
	}
	return &operator{cli: cli}, nil
}

// Scrape returns the content of a listing endpoint on Reddit at path. kind
// should specify the expected return type. See the Kind enum for
// supported options.
func (o *operator) Scrape(
	path,
	after,
	before string,
	limit uint,
	kind Kind,
) ([]Thing, error) {
	req := http.Request{
		Method:     "GET",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Close:      true,
		URL: &url.URL{
			Scheme: "https",
			Host:   oauth2Host,
			Path:   path,
			RawQuery: url.Values{
				"limit":  []string{strconv.Itoa(int(limit))},
				"before": []string{before},
				"after":  []string{after},
			}.Encode(),
		},
		Host: oauth2Host,
	}

	response, err := o.cli.Do(&req)
	if err != nil {
		return nil, err
	}

	return thingListing(response, kind)
}

// GetThing returns a specific thing, by its id. It will return nil if the Thing
// did not exist.
func (o *operator) GetThing(id string, kind Kind) (Thing, error) {
	req := http.Request{
		Method:     "GET",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Close:      true,
		URL: &url.URL{
			Scheme: "https",
			Host:   oauth2Host,
			Path:   "/api/info.json",
			RawQuery: url.Values{
				"id": []string{id},
			}.Encode(),
		},
		Host: oauth2Host,
	}

	response, err := o.cli.Do(&req)
	if err != nil {
		return nil, err
	}

	things, err := thingListing(response, kind)
	if err != nil {
		return nil, err
	}

	if len(things) != 1 {
		return nil, nil
	}

	return things[0], nil
}

// Thread returns a link; the Comments field will be filled with the comment
// tree. Browse each comment's reply tree from the ReplyTree field.
func (o *operator) Thread(permalink string) (*redditproto.Link, error) {
	req := http.Request{
		Method:     "GET",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Close:      true,
		URL: &url.URL{
			Scheme: "https",
			Host:   oauth2Host,
			Path:   fmt.Sprintf("%s.json", permalink),
		},
		Host: oauth2Host,
	}

	response, err := o.cli.Do(&req)
	if err != nil {
		return nil, err
	}

	return parseThread(response)
}

// Inbox returns unread inbox items.
func (o *operator) Inbox() ([]*redditproto.Message, error) {
	req := http.Request{
		Method:     "GET",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Close:      true,
		URL: &url.URL{
			Scheme: "https",
			Host:   oauth2Host,
			Path:   "/message/unread",
		},
		Host: oauth2Host,
	}

	response, err := o.cli.Do(&req)
	if err != nil {
		return nil, err
	}

	return parseInbox(response)
}

// MarkAsRead marks inbox items as read, so they are no longer returned by calls
// to Inbox().
func (o *operator) MarkAsRead(fullnames ...string) error {
	req := http.Request{
		Method:     "POST",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Close:      true,
		URL: &url.URL{
			Scheme: "https",
			Host:   oauth2Host,
			Path:   "/api/read_message",
		},
		Header: formEncoding,
		Body: ioutil.NopCloser(
			bytes.NewBufferString(
				url.Values{
					"id": []string{
						strings.Join(fullnames, ","),
					},
				}.Encode(),
			),
		),
		Host: oauth2Host,
	}

	_, err := o.cli.Do(&req)
	return err
}

// Reply replies to a post, message, or comment.
func (o *operator) Reply(parent, content string) error {
	req := http.Request{
		Method:     "POST",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Close:      true,
		URL: &url.URL{
			Scheme: "https",
			Host:   oauth2Host,
			Path:   "/api/comment",
		},
		Header: formEncoding,
		Body: ioutil.NopCloser(
			bytes.NewBufferString(
				url.Values{
					"thing_id": []string{parent},
					"text":     []string{content},
				}.Encode(),
			),
		),
		Host: oauth2Host,
	}

	_, err := o.cli.Do(&req)
	return err
}

// Compose sends a private message to a user.
func (o *operator) Compose(user, subject, content string) error {
	req := http.Request{
		Method:     "POST",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Close:      true,
		URL: &url.URL{
			Scheme: "https",
			Host:   oauth2Host,
			Path:   "/api/compose",
		},
		Header: formEncoding,
		Body: ioutil.NopCloser(
			bytes.NewBufferString(
				url.Values{
					"to":      []string{user},
					"subject": []string{subject},
					"text":    []string{content},
				}.Encode(),
			),
		),
		Host: oauth2Host,
	}

	_, err := o.cli.Do(&req)
	return err
}

// Submit submits a post.
func (o *operator) Submit(subreddit, kind, title, content string) error {
	req := http.Request{
		Method:     "POST",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Close:      true,
		URL: &url.URL{
			Scheme: "https",
			Host:   oauth2Host,
			Path:   "/api/submit",
		},
		Header: formEncoding,
		Body: ioutil.NopCloser(
			bytes.NewBufferString(
				url.Values{
					"sr":    []string{subreddit},
					"kind":  []string{kind},
					"title": []string{title},
					"url":   []string{content},
					"text":  []string{content},
				}.Encode(),
			),
		),
		Host: oauth2Host,
	}

	_, err := o.cli.Do(&req)
	return err
}

// thingListing returns a slice of Thing interfaces, each referencing a Reddit
// type from a returned JSON listing.
func thingListing(listing io.ReadCloser, kind Kind) ([]Thing, error) {
	if kind == Link {
		links, err := parseLinkListing(listing)
		things := make([]Thing, len(links))
		for i, link := range links {
			things[i] = link
		}
		return things, err
	}

	comments, err := parseCommentListing(listing)
	things := make([]Thing, len(comments))
	for i, link := range comments {
		things[i] = link
	}
	return things, err
}
