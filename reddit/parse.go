package reddit

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

const (
	listingKind = "Listing"
	postKind    = "t3"
	commentKind = "t1"
	messageKind = "t4"
)

// author fields and body fields are set to the deletedKey if the user deletes
// their post.
const deletedKey = "[deleted]"

// thing is a Reddit type that holds all of their subtypes.
type thing struct {
	Kind string                 `json:"kind"`
	Data map[string]interface{} `json:"data"`
}

type listing struct {
	Children []thing `json:"children,omitempty"`
}

// comment wraps the user facing Comment type with a Replies field for
// intermediate parsing.
type comment struct {
	Comment `mapstructure:",squash"`
	Replies thing `mapstructure:"replies"`
}

// parser parses Reddit responses..
type parser interface {
	// parse parses any Reddit response and provides the elements in it.
	parse(blob json.RawMessage) ([]*Comment, []*Post, []*Message, error)
}

type parserImpl struct{}

func newParser() parser {
	return &parserImpl{}
}

// parse parses any Reddit response and provides the elements in it.
func (p *parserImpl) parse(
	blob json.RawMessage,
) ([]*Comment, []*Post, []*Message, error) {
	comments, posts, msgs, listingErr := parseRawListing(blob)
	if listingErr == nil {
		return comments, posts, msgs, nil
	}

	post, threadErr := parseThread(blob)
	if threadErr == nil {
		return nil, []*Post{post}, nil, nil
	}

	return nil, nil, nil, fmt.Errorf(
		"failed to parse as listing [%v] or thread [%v]",
		listingErr, threadErr,
	)
}

// parseRawListing parses a listing json blob and returns the elements in it.
func parseRawListing(
	blob json.RawMessage,
) ([]*Comment, []*Post, []*Message, error) {
	var activityListing thing
	if err := json.Unmarshal(blob, &activityListing); err != nil {
		return nil, nil, nil, err
	}

	return parseListing(&activityListing)
}

// parseThread parses a post from a thread json blob returned by Reddit.
//
// Reddit structures this as two things in an array, the first thing being a
// listing with only the post and the second thing being a listing of comments.
func parseThread(blob json.RawMessage) (*Post, error) {
	var listings [2]thing
	if err := json.Unmarshal(blob, &listings); err != nil {
		return nil, err
	}

	_, posts, _, err := parseListing(&listings[0])
	if err != nil {
		return nil, err
	}

	if len(posts) != 1 {
		return nil, fmt.Errorf("expected 1 post; found %d", len(posts))
	}

	comments, _, _, err := parseListing(&listings[1])
	if err != nil {
		return nil, err
	}

	posts[0].Replies = comments
	return posts[0], nil
}

// parseListing parses a Reddit listing type and returns the elements inside it.
func parseListing(t *thing) ([]*Comment, []*Post, []*Message, error) {
	if t.Kind != listingKind {
		return nil, nil, nil, fmt.Errorf("thing is not listing")
	}

	l := &listing{}
	if err := mapstructure.Decode(t.Data, l); err != nil {
		return nil, nil, nil, mapDecodeError(err, t.Data)
	}

	comments := []*Comment{}
	posts := []*Post{}
	msgs := []*Message{}
	err := error(nil)

	for _, c := range l.Children {
		if err != nil {
			break
		}

		var comment *Comment
		var post *Post
		var msg *Message

		// Reddit sets the "Kind" field of comments in the inbox, which
		// have only Message and not Comment fields, to commentKind. The
		// give away in this case is that comments in message form have
		// a field called "was_comment". Reddit does this because they
		// hate programmers.
		if c.Kind == messageKind || c.Data["was_comment"] != nil {
			msg, err = parseMessage(&c)
			msgs = append(msgs, msg)
		} else if c.Kind == commentKind {
			comment, err = parseComment(&c)
			comments = append(comments, comment)
		} else if c.Kind == postKind {
			post, err = parsePost(&c)
			posts = append(posts, post)
		}
	}

	return comments, posts, msgs, err
}

// parseComment parses a comment into the user facing Comment struct.
func parseComment(t *thing) (*Comment, error) {
	// Reddit makes the replies field a string if it is empty, just to make
	// it harder for programmers who like static type systems.
	value, present := t.Data["replies"]
	if present {
		if str, ok := value.(string); ok && str == "" {
			delete(t.Data, "replies")
		}
	}

	c := &comment{}
	if err := mapstructure.Decode(t.Data, c); err != nil {
		return nil, mapDecodeError(err, t.Data)
	}

	var err error
	if c.Replies.Kind == listingKind {
		c.Comment.Replies, _, _, err = parseListing(&c.Replies)
	}

	c.Comment.Deleted = c.Comment.Body == deletedKey
	return &c.Comment, err
}

// parsePost parses a post into the user facing Post struct.
func parsePost(t *thing) (*Post, error) {
	p := &Post{}
	if err := mapstructure.Decode(t.Data, p); err != nil {
		return nil, mapDecodeError(err, t.Data)
	}

	p.Deleted = p.SelfText == deletedKey
	return p, nil
}

// parseMessage parses a message into the user facing Message struct.
func parseMessage(t *thing) (*Message, error) {
	m := &Message{}
	return m, mapstructure.Decode(t.Data, m)
}

func mapDecodeError(err error, val interface{}) error {
	return fmt.Errorf(
		"failed to decode json map into struct: %v; value: %v",
		err, val,
	)
}
