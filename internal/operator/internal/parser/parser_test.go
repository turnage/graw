package parser

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/turnage/redditproto"
)

// rcloser wraps a buffer so that it can passed as an io.ReadCloser.
type rcloser struct {
	*bytes.Buffer
}

func (b *rcloser) Close() error {
	return nil
}

type errReader struct{}

func (e *errReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("an error")
}

func (e *errReader) Close() error {
	return nil
}

func TestParseLinkListing(t *testing.T) {
	if _, err := ParseLinkListing(nil); err == nil {
		t.Errorf("wanted error for nil content")
	}

	if _, err := ParseLinkListing(
		&rcloser{bytes.NewBufferString(`[]"`)},
	); err == nil {
		t.Errorf("wanted error for invalid json")
	}

	if _, err := ParseLinkListing(
		&rcloser{bytes.NewBufferString(`{}`)},
	); err == nil {
		t.Errorf("wanted error for nil data")
	}

	if _, err := ParseLinkListing(
		&rcloser{bytes.NewBufferString(`{"data": {}}`)},
	); err == nil {
		t.Errorf("wanted error for nil children")
	}

	links, err := ParseLinkListing(&rcloser{bytes.NewBufferString(`{
		"data": {
			"children": [
				{"data": {"title": "hello"}},
				{"data": {"title": "hola"}},
				{"data": {"title": "bye"}}
			]
		}
	}`)})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(links) != 3 {
		t.Errorf("wanted to find 3 links; resp is %v", links)
	}
}

func TestParseThread(t *testing.T) {
	if _, err := ParseThread(nil); err == nil {
		t.Errorf("wanted error for nil content")
	}

	if _, err := ParseThread(&errReader{}); err == nil {
		t.Errorf("Wanted error for failed read")
	}

	if _, err := ParseThread(
		&rcloser{bytes.NewBufferString(`asds`)},
	); err == nil {
		t.Errorf("wanted error for invalid json")
	}

	if _, err := ParseThread(
		&rcloser{bytes.NewBufferString(`[]`)},
	); err == nil {
		t.Errorf("wanted error for missing listings")
	}

	if _, err := ParseThread(&rcloser{bytes.NewBufferString(`[
			{"data": {}},
			{"data": {}}
	]`)}); err == nil {
		t.Errorf("wanted error for bad link listing")
	}

	if _, err := ParseThread(&rcloser{bytes.NewBufferString(`[
		{
			"data": {
				"children": [
					{"data": {"title": "hola"}},
					{"data": {"title": "bye"}}
				]
			}
		},
		{
			"data": {
				"children": [
					{"data": {"id": "arnold"}},
					{"data": {"id": "harold"}}
				]
			}
		}
	]`)}); err == nil {
		t.Errorf("wanted error for non-one link listing")
	}

	thread, err := ParseThread(&rcloser{bytes.NewBufferString(`[
		{
			"data": {
				"children": [
					{"data": {"title": "hola"}}
				]
			}
		},
		{
			"data": {
				"children": [
					{"data": {"id": "arnold"}},
					{"data": {"id": "harold"}}
				]
			}
		}
	]`)})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if thread.GetTitle() != "hola" {
		t.Errorf("got %s; wanted hola", thread.GetTitle())
	}
	if len(thread.GetComments()) != 2 {
		t.Errorf("got %v; wanted 2 comments", thread.GetComments)
	}

	if _, err := ParseThread(&rcloser{bytes.NewBufferString(`[
		{
			"data": {
				"children": [
					{"data": {"title": "hola"}}
				]
			}
		},
		{
			"data": {
				"children": [
					{"data": {"id": "arnold"}},
					{"data": {
							"replies": "",
							"id": "harold"
						}
					}
				]
			}
		}
	]`)}); err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestParseInbox(t *testing.T) {
	if _, err := ParseInbox(nil); err == nil {
		t.Errorf("wanted error for nil content")
	}

	if _, err := ParseInbox(
		&rcloser{bytes.NewBufferString(`[]"`)},
	); err == nil {
		t.Errorf("wanted error for invalid json")
	}

	if _, err := ParseInbox(
		&rcloser{bytes.NewBufferString(`{}`)},
	); err == nil {
		t.Errorf("wanted error for nil data")
	}

	if _, err := ParseInbox(
		&rcloser{bytes.NewBufferString(`{"data": {}}`)},
	); err == nil {
		t.Errorf("wanted error for nil children")
	}

	messages, err := ParseInbox(&rcloser{bytes.NewBufferString(`{
		"data": {
			"children" : [
				{"data": {"was_comment": true}}
			]
		}
	}`)})
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("got %d messages; wanted 1", len(messages))
	}

	if !messages[0].GetWasComment() {
		t.Errorf("wanted message to be a comment")
	}
}

func TestUnpackCommentListing(t *testing.T) {
	if unpackCommentListing(&redditproto.CommentListing{}) != nil {
		t.Errorf("wanted empty slice when data is nil")
	}

	if unpackCommentListing(&redditproto.CommentListing{
		Data: &redditproto.CommentData{},
	}) != nil {
		t.Errorf("wanted empty slice when children is nil")
	}

	topId := "top"
	subId := "sub"
	sublisting := &redditproto.CommentListing{
		Data: &redditproto.CommentData{
			Children: []*redditproto.CommentChildren{
				&redditproto.CommentChildren{
					Data: &redditproto.Comment{
						Id: &subId,
					},
				},
			},
		},
	}
	listing := &redditproto.CommentListing{
		Data: &redditproto.CommentData{
			Children: []*redditproto.CommentChildren{
				&redditproto.CommentChildren{
					Data: &redditproto.Comment{
						Id:      &topId,
						Replies: sublisting,
					},
				},
			},
		},
	}
	expected := []*redditproto.Comment{
		&redditproto.Comment{
			Id: &topId,
			ReplyTree: []*redditproto.Comment{
				&redditproto.Comment{
					Id: &subId,
				},
			},
		},
	}
	actual := unpackCommentListing(listing)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v; wanted %v", actual, expected)
	}
}
