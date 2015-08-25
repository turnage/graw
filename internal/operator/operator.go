// Package operator makes api calls to Reddit.
package operator

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/turnage/graw/internal/operator/internal/client"
	"github.com/turnage/graw/internal/operator/internal/parser"
	"github.com/turnage/graw/internal/operator/internal/request"
	"github.com/turnage/redditproto"
)

// Operator makes api calls to Reddit.
type Operator struct {
	cli client.Client
}

// New returns a new operator which uses agent as its identity. agent should be
// a filename of a file containing a UserAgent protobuffer.
func New(agent string) (*Operator, error) {
	cli, err := client.New(agent)
	if err != nil {
		return nil, err
	}
	return &Operator{cli: cli}, nil
}

// NewMock returns an operator which will act as if it receives this provided
// response from the server for all requests.
func NewMock(response string) *Operator {
	return &Operator{cli: client.NewMock(response)}
}

// Scrape returns posts from a subreddit, in the specified sort order, with the
// specified reference points for direction, up to lim. lims above 100 are
// ineffective because Reddit will return only 100 posts per query. The Comments
// field will not be filled. For comments, request a thread using Thread().
func (o *Operator) Scrape(sub, sort, after, before string, lim uint) ([]*redditproto.Link, error) {
	req, err := request.New(
		"GET",
		fmt.Sprintf("https://oauth.reddit.com/r/%s/%s.json", sub, sort),
		&url.Values{
			"limit":  []string{strconv.Itoa(int(lim))},
			"after":  []string{after},
			"before": []string{before},
		},
	)
	if err != nil {
		return nil, err
	}

	response, err := o.cli.Do(req)
	if err != nil {
		return nil, err
	}

	return parser.ParseLinkListing(response)
}

// Threads returns specific threads, requested by their fullname (t3_[id]).
// The Comments field will be not be filled. For comments, request a thread
// using Thread().
func (o *Operator) Threads(fullnames ...string) ([]*redditproto.Link, error) {
	ids := strings.Join(fullnames, ",")
	req, err := request.New(
		"GET",
		fmt.Sprintf("https://oauth.reddit.com/by_id/%s", ids),
		nil,
	)
	if err != nil {
		return nil, err
	}

	response, err := o.cli.Do(req)
	if err != nil {
		return nil, err
	}

	return parser.ParseLinkListing(response)
}

// Thread returns a link; the Comments field will be filled with the comment
// tree. Browse each comment's reply tree from the ReplyTree field.
func (o *Operator) Thread(permalink string) (*redditproto.Link, error) {
	req, err := request.New(
		"GET",
		fmt.Sprintf("https://oauth.reddit.com%s.json", permalink),
		nil,
	)
	if err != nil {
		return nil, err
	}

	response, err := o.cli.Do(req)
	if err != nil {
		return nil, err
	}

	return parser.ParseThread(response)
}

// Inbox returns unread messages and marks them as read.
func (o *Operator) Inbox() ([]*redditproto.Message, error) {
	req, err := request.New(
		"GET",
		"https://oauth.reddit.com/message/unread",
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

	messages, err := parser.ParseInbox(response)
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
		"https://oauth.reddit.com/api/read_message",
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
func (o *Operator) Reply(parent, content string) error {
	req, err := request.New(
		"POST",
		"https://oauth.reddit.com/api/comment",
		&url.Values{
			"thing_id": []string{parent},
			"text":     []string{content},
		},
	)
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
func (o *Operator) Compose(user, subject, content string) error {
	req, err := request.New(
		"POST",
		"https://oauth.reddit.com/api/compose",
		&url.Values{
			"to":      []string{user},
			"subject": []string{subject},
			"text":    []string{content},
		},
	)
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
func (o *Operator) Submit(subreddit, kind, title, content string) error {
	params := &url.Values{
		"sr":    []string{subreddit},
		"kind":  []string{kind},
		"title": []string{title},
	}
	if kind == "link" {
		params.Add("url", content)
	} else if kind == "self" {
		params.Add("text", content)
	} else {
		return fmt.Errorf("unsupported post type")
	}

	req, err := request.New(
		"POST",
		"https://oauth.reddit.com/api/submit",
		params,
	)
	if err != nil {
		return err
	}

	_, err = o.cli.Do(req)
	if err != nil {
		return err
	}

	return nil
}
