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

const (
	// MaxLinks is the amount of posts reddit will return for a scrape
	// query.
	MaxLinks = 100
	// baseURL is the url all requests extend from.
	baseURL = "https://oauth.reddit.com"
	// oauth2Host is the hostname of Reddit's OAuth2 server.
	oauth2Host = "oauth.reddit.com"
	// deletedAuthor is the author value if a post or comment was deleted.
	deletedAuthor = "[deleted]"
)

var (
	formEncoding = map[string][]string{
		"content-type": {"application/x-www-form-urlencoded"},
	}
)

// Operator makes api calls to Reddit.
type Operator interface {
	// Posts returns a listing of posts in a subreddit.
	Posts(subreddit, after, before string, limit uint) ([]*redditproto.Link, error)
	// UserContent returns a listing of posts and comments from a user's
	// landing page.
	UserContent(user, after, before string, limit uint) ([]*redditproto.Link, []*redditproto.Comment, error)
	// IsThereThing fetches a particular thing from reddit. Thing returns
	// whether there is such a thing.
	IsThereThing(id string) (bool, error)
	// Thread fetches a post and its comment tree.
	Thread(permalink string) (*redditproto.Link, error)
	// Inbox fetches unread messages from the reddit inbox.
	Inbox() ([]*redditproto.Message, error)
	// MarkAsRead marks inbox items read.
	MarkAsRead(fullnames ...string) error
	// Reply replies to reddit item.
	Reply(parent, content string) error
	// Compose sends a private message to a user.
	Compose(user, subject, content string) error
	// Submit posts to Reddit.
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

// Posts returns a list of posts in a subreddit.
func (o *operator) Posts(
	subreddit,
	after,
	before string,
	limit uint,
) ([]*redditproto.Link, error) {
	bytes, err := o.exec(
		http.Request{
			Method:     "GET",
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Close:      true,
			URL: &url.URL{
				Scheme:   "https",
				Host:     oauth2Host,
				Path:     fmt.Sprintf("/r/%s/new", subreddit),
				RawQuery: listingParams(limit, after, before),
			},
			Host: oauth2Host,
		},
	)
	if err != nil {
		return nil, err
	}

	return redditproto.ParseLinkListing(bytes)
}

// UserContent returns a list of content from a user's landing page.
func (o *operator) UserContent(
	user,
	after,
	before string,
	limit uint,
) ([]*redditproto.Link, []*redditproto.Comment, error) {
	bytes, err := o.exec(
		http.Request{
			Method:     "GET",
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Close:      true,
			URL: &url.URL{
				Scheme:   "https",
				Host:     oauth2Host,
				Path:     fmt.Sprintf("/u/%s.json", user),
				RawQuery: listingParams(limit, after, before),
			},
			Host: oauth2Host,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	return redditproto.ParseComboListing(bytes)
}

// IsThereThing returns whether a thing by the given id exists.
func (o *operator) IsThereThing(id string) (bool, error) {
	bytes, err := o.exec(
		http.Request{
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
		},
	)

	if err != nil {
		return false, err
	}

	links, comments, err := redditproto.ParseComboListing(bytes)
	if err != nil {
		return false, err
	}

	if len(links) == 1 {
		return links[0].GetAuthor() != deletedAuthor, nil
	}

	if len(comments) == 1 {
		return comments[0].GetAuthor() != deletedAuthor, nil
	}

	return false, nil
}

// Thread returns a link; the Comments field will be filled with the comment
// tree. Browse each comment's reply tree from the ReplyTree field.
func (o *operator) Thread(permalink string) (*redditproto.Link, error) {
	bytes, err := o.exec(
		http.Request{
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
		},
	)
	if err != nil {
		return nil, err
	}

	return redditproto.ParseThread(bytes)
}

// Inbox returns unread inbox items.
func (o *operator) Inbox() ([]*redditproto.Message, error) {
	bytes, err := o.exec(
		http.Request{
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
		},
	)
	if err != nil {
		return nil, err
	}

	return redditproto.ParseMessageListing(bytes)
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

// exec executes a request and returns the response body bytes.
func (o *operator) exec(r http.Request) ([]byte, error) {
	response, err := o.cli.Do(&r)
	if err != nil {
		return nil, err
	}

	return responseBytes(response)
}

// listingParams returns encoded values for parameters to a Reddit listing
// endpoint.
func listingParams(limit uint, after, before string) string {
	return url.Values{
		"limit":  []string{strconv.Itoa(int(limit))},
		"before": []string{before},
		"after":  []string{after},
	}.Encode()
}

// responseBytes returns a slice of bytes from a response body.
func responseBytes(response io.ReadCloser) ([]byte, error) {
	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(response); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
