package graw

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
	ID   string                 `json:"id,omitempty"`
	Name string                 `json:"name,omitempty"`
	Kind string                 `json:"kind"`
	Data map[string]interface{} `json:"data"`
}

type listing struct {
	Children []thing `json:"children,omitempty"`
}

// comment wraps the user facing Comment data type with a Replies field for
// intermediate parsing.
type comment struct {
	Comment `mapstructure:",squash"`
	Replies thing `mapstructure:"replies"`
}

func thingKindError(actual, wanted string) error {
	return fmt.Errorf("got kind %s; wanted kind %s", actual, wanted)
}

func mapDecodeError(err error, val interface{}) error {
	return fmt.Errorf(
		"failed to decode json map into struct: %v; value: %v",
		err, val,
	)
}

// propagateThingMetadataDown takes fields from thing and adds them to the child
// map so the child knows when it was created and what its full name is.
func propagateThingMetadataDown(t *thing) {
	t.Data["id"] = t.ID
	t.Data["name"] = t.Name
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

	_, posts, err := parseListing(&listings[0])
	if err != nil {
		return nil, err
	}

	comments, _, err := parseListing(&listings[1])
	if err != nil {
		return nil, err
	}

	if len(posts) != 1 {
		return nil, fmt.Errorf("expected 1 post; found %d", len(posts))
	}

	posts[0].Replies = comments
	return posts[0], nil
}

// parseListing parses a Reddit listing type and returns the elements inside it.
func parseListing(t *thing) ([]*Comment, []*Post, error) {
	if t.Kind != listingKind {
		return nil, nil, thingKindError(t.Kind, listingKind)
	}

	l := &listing{}
	if err := mapstructure.Decode(t.Data, l); err != nil {
		return nil, nil, mapDecodeError(err, t.Data)
	}

	comments := []*Comment{}
	posts := []*Post{}
	err := error(nil)

	for _, c := range l.Children {
		if err != nil {
			break
		}

		var comment *Comment
		var post *Post
		switch c.Kind {
		case commentKind:
			comment, err = parseComment(&c)
			comments = append(comments, comment)
		case postKind:
			post, err = parsePost(&c)
			posts = append(posts, post)
		}
	}

	return comments, posts, err
}

// parseComment parses a comment into the user facing Comment struct.
func parseComment(t *thing) (*Comment, error) {
	if t.Kind != commentKind {
		return nil, thingKindError(t.Kind, commentKind)
	}

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
		c.Comment.Replies, _, err = parseListing(&c.Replies)
	}

	c.Comment.Deleted = c.Comment.Body == deletedKey
	return &c.Comment, err
}

// parsePost parses a post into the user facing Post struct.
func parsePost(t *thing) (*Post, error) {
	if t.Kind != postKind {
		return nil, thingKindError(t.Kind, postKind)
	}

	p := &Post{}
	if err := mapstructure.Decode(t.Data, p); err != nil {
		return nil, mapDecodeError(err, t.Data)
	}

	p.Deleted = p.SelfText == deletedKey
	return p, nil
}
