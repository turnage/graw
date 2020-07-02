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
	moreKind    = "more"
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

type more struct {
	Errors []interface{} `json:"errors,omitempty"`
	Data   []thing       `json:"data"`
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
	parse(blob json.RawMessage) ([]*Comment, []*Post, []*Message, []*More, error)
	parse_submitted(blob json.RawMessage) (Submission, error)
}

type parserImpl struct{}

func newParser() parser {
	return &parserImpl{}
}

// parse parses any Reddit response and provides the elements in it.
func (p *parserImpl) parse(
	blob json.RawMessage,
) ([]*Comment, []*Post, []*Message, []*More, error) {
	comments, posts, msgs, mores, listingErr := parseRawListing(blob)
	if listingErr == nil {
		return comments, posts, msgs, mores, nil
	}

	post, mores, threadErr := parseThread(blob)
	if threadErr == nil {
		return nil, []*Post{post}, nil, mores, nil
	}

	comments, posts, msgs, mores, moreErr := parseMoreChildren(blob)
	if moreErr == nil {
		return comments, posts, msgs, mores, nil
	}

	return nil, nil, nil, nil, fmt.Errorf(
		"failed to parse as listing [%v], thread [%v], or more [%v]",
		listingErr, threadErr,
	)
}

// parse_submitted parses a response from reddit describing
// the status of some resource that was submitted
func (p *parserImpl) parse_submitted(blob json.RawMessage) (Submission, error) {
	var wrapped map[string]interface{}
	err := json.Unmarshal(blob, &wrapped)
	if err != nil {
		return Submission{}, err
	}

	wrapped = wrapped["json"].(map[string]interface{})
	if len(wrapped["errors"].([]interface{})) != 0 {
		return Submission{}, fmt.Errorf("API errors were returned: %v", wrapped["errors"])
	}

	data := wrapped["data"].(map[string]interface{})

	// Comment submissions are further wrapped in a things block
	// because of ... something? There only appears to be a single thing
	// This transformes var data to be the data of the single thing
	// This also mirrors https://reddit.com/ + permalink -> url
	things, has_things := data["things"].([]interface{})
	if has_things && len(things) == 1 {
		data = things[0].(map[string]interface{})["data"].(map[string]interface{})
		data["url"] = fmt.Sprintf("https://reddit.com%s", data["permalink"])
	}

	var submission Submission
	err = mapstructure.Decode(data, &submission)
	return submission, err
}

// parseRawListing parses a listing json blob and returns the elements in it.
func parseRawListing(
	blob json.RawMessage,
) ([]*Comment, []*Post, []*Message, []*More, error) {
	var activityListing thing
	if err := json.Unmarshal(blob, &activityListing); err != nil {
		return nil, nil, nil, nil, err
	}

	return parseListing(&activityListing)
}

// parseMoreChildren parses the json blob from morechildren calls and returns the elements in it.
func parseMoreChildren(
	blob json.RawMessage,
) ([]*Comment, []*Post, []*Message, []*More, error) {
	var wrapped map[string]interface{}
	err := json.Unmarshal(blob, &wrapped)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	wrapped = wrapped["json"].(map[string]interface{})
	if len(wrapped["errors"].([]interface{})) != 0 {
		return nil, nil, nil, nil, fmt.Errorf("API errors were returned: %v", wrapped["errors"])
	}

	data := wrapped["data"].(map[string]interface{})
	// More submissions are further wrapped in a things block,
	// so reorganize data so that it makes sense
	things, hasThings := data["things"].([]interface{})
	if hasThings {
		data["data"] = things
		delete(data, "things")
	} else {
		return nil, nil, nil, nil, fmt.Errorf("No thing types returned")
	}

	var m more
	err = mapstructure.Decode(data, &m)

	if err != nil {
		return nil, nil, nil, nil, err
	} else if m.Errors != nil {
		return nil, nil, nil, nil, fmt.Errorf("%v", m.Errors)
	}

	return parseChildren(m.Data)
}

// parseThread parses a post from a thread json blob returned by Reddit.
//
// Reddit structures this as two things in an array, the first thing being a
// listing with only the post and the second thing being a listing of comments.
func parseThread(blob json.RawMessage) (*Post, []*More, error) {
	var listings [2]thing
	if err := json.Unmarshal(blob, &listings); err != nil {
		return nil, nil, err
	}

	_, posts, _, _, err := parseListing(&listings[0])
	if err != nil {
		return nil, nil, err
	}

	if len(posts) != 1 {
		return nil, nil, fmt.Errorf("expected 1 post; found %d", len(posts))
	}

	comments, _, _, mores, err := parseListing(&listings[1])
	if err != nil {
		return nil, nil, err
	}

	posts[0].Replies = comments
	return posts[0], mores, nil
}

// parseListing parses a Reddit listing type and returns the elements inside it.
func parseListing(t *thing) ([]*Comment, []*Post, []*Message, []*More, error) {
	if t.Kind != listingKind {
		return nil, nil, nil, nil, fmt.Errorf("thing is not listing")
	}

	l := &listing{}
	if err := mapstructure.Decode(t.Data, l); err != nil {
		return nil, nil, nil, nil, mapDecodeError(err, t.Data)
	}

	return parseChildren(l.Children)
}

// parseChildren returns a list of parsed objects from the given list of things
func parseChildren(children []thing) ([]*Comment, []*Post, []*Message, []*More, error) {
	comments := []*Comment{}
	posts := []*Post{}
	msgs := []*Message{}
	mores := []*More{}
	err := error(nil)

	for _, c := range children {
		if err != nil {
			break
		}

		var comment *Comment
		var post *Post
		var msg *Message
		var more *More

		// Reddit sets the "Kind" field of comments in the inbox, which
		// have only Message and not Comment fields, to commentKind. The
		// give away in this case is that comments in message form have
		// a field called "was_comment". Reddit does this because they
		// hate programmers.
		if c.Kind == messageKind || c.Data["was_comment"] != nil {
			msg, err = parseMessage(&c)
			msgs = append(msgs, msg)
		} else if c.Kind == commentKind {
			comment, more, err = parseComment(&c)
			comments = append(comments, comment)
			if more != nil {
				mores = append(mores, more)
			}
		} else if c.Kind == postKind {
			post, err = parsePost(&c)
			posts = append(posts, post)
		} else if c.Kind == moreKind {
			more, err = parseMore(&c)
			mores = append(mores, more)
		}
	}

	return comments, posts, msgs, mores, err
}

// parseComment parses a comment into the user facing Comment struct.
func parseComment(t *thing) (*Comment, *More, error) {
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
		return nil, nil, mapDecodeError(err, t.Data)
	}

	var err error
	var mores []*More
	if c.Replies.Kind == listingKind {
		c.Comment.Replies, _, _, mores, err = parseListing(&c.Replies)
	}

	c.Comment.Deleted = c.Comment.Body == deletedKey
	// we should only evey have one more values per comment branch
	if len(mores) == 1 {
		return &c.Comment, mores[0], err
	}
	return &c.Comment, nil, err
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

// parseMore parses a more comment list into the user facing More struct.
func parseMore(t *thing) (*More, error) {
	m := &More{}
	if err := mapstructure.Decode(t.Data, m); err != nil {
		return nil, mapDecodeError(err, t.Data)
	}

	return m, nil
}

func mapDecodeError(err error, val interface{}) error {
	return fmt.Errorf(
		"failed to decode json map into struct: %v; value: %v",
		err, val,
	)
}
