// Package operator makes api calls to Reddit.
package operator

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/turnage/graw/internal/operator/internal/client"
	"github.com/turnage/graw/internal/operator/internal/request"
	"github.com/turnage/redditproto"
)

const (
	// MaxLinks is the amount of posts reddit will return for a scrape
	// query.
	MaxLinks = 100
	// baseURL is the url all requests extend from.
	baseURL = "https://oauth.reddit.com"
)

// Operator makes api calls to Reddit.
type Operator interface {
	// Scrape fetches new reddit posts (see definition).
	Scrape(subreddit, sort, after, before string, limit uint) ([]*redditproto.Link, error)
	// Threads fetches specific threads by name (see definition).
	Threads(fullnames ...string) ([]*redditproto.Link, error)
	// Thread fetches a post and its comment tree (see definition).
	Thread(permalink string) (*redditproto.Link, error)
	// Inbox fetches unread messages from the reddit inbox (see definition).
	Inbox() ([]*redditproto.Message, error)
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

// Scrape returns posts from a subreddit, in the specified sort order, with the
// specified reference points for direction, up to limit. The Comments
// field will not be filled. For comments, request a thread using Thread().
func (o *operator) Scrape(
	subreddit,
	sort,
	after,
	before string,
	limit uint,
) ([]*redditproto.Link, error) {
	req, err := scrapeRequest(subreddit, sort, after, before, limit)
	if err != nil {
		return nil, err
	}

	response, err := o.cli.Do(req)
	if err != nil {
		return nil, err
	}

	return parseLinkListing(response)
}

// Threads returns specific threads, requested by their fullname (t3_[id]).
// The Comments field will be not be filled. For comments, request a thread
// using Thread().
func (o *operator) Threads(fullnames ...string) ([]*redditproto.Link, error) {
	req, err := threadsRequest(fullnames)
	if err != nil {
		return nil, err
	}

	response, err := o.cli.Do(req)
	if err != nil {
		return nil, err
	}

	return parseLinkListing(response)
}

// Thread returns a link; the Comments field will be filled with the comment
// tree. Browse each comment's reply tree from the ReplyTree field.
func (o *operator) Thread(permalink string) (*redditproto.Link, error) {
	req, err := threadRequest(permalink)
	if err != nil {
		return nil, err
	}

	response, err := o.cli.Do(req)
	if err != nil {
		return nil, err
	}

	return parseThread(response)
}

// Inbox returns unread messages and marks them as read.
func (o *operator) Inbox() ([]*redditproto.Message, error) {
	req, err := request.New(
		"GET",
		fmt.Sprintf("%s/message/unread", baseURL),
		&url.Values{
			"limit": []string{"100"},
		},
	)
	if err != nil {
		return nil, err
	}

	response, err := o.cli.Do(req)
	if err != nil {
		return nil, err
	}

	messages, err := parseInbox(response)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return messages, nil
	}

	messageIds := make([]string, len(messages))
	for i, message := range messages {
		messageIds[i] = message.GetName()
	}

	req, err = request.New(
		"POST",
		fmt.Sprintf("%s/api/read_message", baseURL),
		&url.Values{
			"id": messageIds,
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = o.cli.Do(req)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// Reply replies to a post, message, or comment.
func (o *operator) Reply(parent, content string) error {
	req, err := replyRequest(parent, content)
	if err != nil {
		return err
	}

	_, err = o.cli.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// Compose sends a private message to a user.
func (o *operator) Compose(user, subject, content string) error {
	req, err := composeRequest(user, subject, content)
	if err != nil {
		return err
	}

	_, err = o.cli.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// Submit submits a post.
func (o *operator) Submit(subreddit, kind, title, content string) error {
	req, err := submitRequest(subreddit, kind, title, content)
	if err != nil {
		return err
	}

	_, err = o.cli.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// scrapeRequests returns an http request representing a subreddit scrape query
// to the Reddit api.
func scrapeRequest(
	subreddit,
	sort,
	before,
	after string,
	limit uint,
) (*http.Request, error) {
	if limit > MaxLinks {
		return nil, fmt.Errorf(
			"%s links requested; max is %s",
			limit,
			MaxLinks)
	}

	if subreddit == "" {
		return nil, fmt.Errorf("no subreddit provided")
	}

	if sort == "" {
		return nil, fmt.Errorf("no sort provided")
	}

	if limit == 0 {
		return nil, fmt.Errorf("no request necessary for 0 links")
	}

	if before != "" && after != "" {
		return nil, fmt.Errorf("have both after and before ids; " +
			"this tells reddit to scrape in two directions.")
	}

	params := &url.Values{
		"limit": []string{strconv.Itoa(int(limit))},
	}
	if before != "" {
		params.Add("before", before)
	} else if after != "" {
		params.Add("after", after)
	}

	return request.New(
		"GET",
		fmt.Sprintf("%s/r/%s/%s.json", baseURL, subreddit, sort),
		params,
	)
}

// threadsRequest creates an http request that represents a by_id api call to
// Reddit.
func threadsRequest(fullnames []string) (*http.Request, error) {
	if len(fullnames) == 0 {
		return nil, fmt.Errorf("no threads provided")
	}

	return request.New(
		"GET",
		fmt.Sprintf(
			"%s/by_id/%s",
			baseURL,
			strings.Join(fullnames, ","),
		),
		nil,
	)
}

// threadRequest creates an http request that represents a call for a specific
// thread comment listing from Reddit.
func threadRequest(permalink string) (*http.Request, error) {
	if permalink == "" {
		return nil, fmt.Errorf("no permalink with which to find thread")
	}

	return request.New(
		"GET",
		fmt.Sprintf("%s%s.json", baseURL, permalink),
		nil,
	)
}

// replyRequest creates an http request that represents a reply api call to
// Reddit.
func replyRequest(parent, content string) (*http.Request, error) {
	if parent == "" {
		return nil, fmt.Errorf("no parent provided to reply to")
	}

	if content == "" {
		return nil, fmt.Errorf("reply body empty")
	}

	return request.New(
		"POST",
		fmt.Sprintf("%s/api/comment", baseURL),
		&url.Values{
			"thing_id": []string{parent},
			"text":     []string{content},
		},
	)
}

// composeRequest creates an http request that represents a compose api call to
// Reddit.
func composeRequest(user, subject, content string) (*http.Request, error) {
	if user == "" {
		return nil, fmt.Errorf("no user to message")
	}

	if subject == "" {
		return nil, fmt.Errorf("no subject for message")
	}

	if content == "" {
		return nil, fmt.Errorf("no body for message")
	}

	return request.New(
		"POST",
		fmt.Sprintf("%s/api/compose", baseURL),
		&url.Values{
			"to":      []string{user},
			"subject": []string{subject},
			"text":    []string{content},
		},
	)
}

// submitRequest creates an http request that represents a submit api call to
// Reddit.
func submitRequest(
	subreddit,
	kind,
	title,
	content string,
) (*http.Request, error) {
	if subreddit == "" {
		return nil, fmt.Errorf("no subreddit provided")
	}

	if title == "" {
		return nil, fmt.Errorf("no title provided")
	}

	params := &url.Values{
		"sr":    []string{subreddit},
		"kind":  []string{kind},
		"title": []string{title},
	}
	if kind == "link" {
		if content == "" {
			return nil, fmt.Errorf("no url provided for link post")
		}
		params.Add("url", content)
	} else if kind == "self" {
		params.Add("text", content)
	} else {
		return nil, fmt.Errorf("unsupported post type")
	}

	return request.New(
		"POST",
		fmt.Sprintf("%s/api/submit", baseURL),
		params,
	)
}
