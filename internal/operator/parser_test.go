package operator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/turnage/redditproto"
)

type errReader struct{}

func (e *errReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("an error")
}

func (e *errReader) Close() error {
	return nil
}

func TestParseLinkListing(t *testing.T) {
	if _, err := parseLinkListing(nil); err == nil {
		t.Errorf("wanted error for nil content")
	}

	if _, err := parseLinkListing(
		ioutil.NopCloser(bytes.NewBufferString(`[]"`)),
	); err == nil {
		t.Errorf("wanted error for invalid json")
	}

	if _, err := parseLinkListing(
		ioutil.NopCloser(bytes.NewBufferString(`{}`)),
	); err == nil {
		t.Errorf("wanted error for nil data")
	}

	if _, err := parseLinkListing(
		ioutil.NopCloser(bytes.NewBufferString(`{"data": {}}`)),
	); err == nil {
		t.Errorf("wanted error for nil children")
	}

	links, err := parseLinkListing(ioutil.NopCloser(bytes.NewBufferString(`{
		"data": {
			"children": [
				{"data": {"title": "hello"}},
				{"data": {"title": "hola"}},
				{"data": {"title": "bye"}}
			]
		}
	}`)))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(links) != 3 {
		t.Errorf("wanted to find 3 links; resp is %v", links)
	}
}

func TestParseCommentListing(t *testing.T) {
	if _, err := parseCommentListing(nil); err == nil {
		t.Errorf("wanted error for nil content")
	}

	if _, err := parseCommentListing(
		ioutil.NopCloser(bytes.NewBufferString(`[]"`)),
	); err == nil {
		t.Errorf("wanted error for invalid json")
	}

	if _, err := parseCommentListing(
		ioutil.NopCloser(bytes.NewBufferString(`{}`)),
	); err == nil {
		t.Errorf("wanted error for nil data")
	}

	if _, err := parseCommentListing(
		ioutil.NopCloser(bytes.NewBufferString(`{"data": {}}`)),
	); err == nil {
		t.Errorf("wanted error for nil children")
	}

	comments, err := parseCommentListing(ioutil.NopCloser(bytes.NewBufferString(`{
		"data": {
			"children": [
				{"data": {"body": "hello"}},
				{"data": {"body": "hola"}},
				{"data": {"body": "bye"}}
			]
		}
	}`)))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(comments) != 3 {
		t.Errorf("wanted to find 3 comments; resp is %v", comments)
	}
}

func TestParseThread(t *testing.T) {
	if _, err := parseThread(nil); err == nil {
		t.Errorf("wanted error for nil content")
	}

	if _, err := parseThread(&errReader{}); err == nil {
		t.Errorf("Wanted error for failed read")
	}

	if _, err := parseThread(
		ioutil.NopCloser(bytes.NewBufferString(`asds`)),
	); err == nil {
		t.Errorf("wanted error for invalid json")
	}

	if _, err := parseThread(
		ioutil.NopCloser(bytes.NewBufferString(`[]`)),
	); err == nil {
		t.Errorf("wanted error for missing listings")
	}

	if _, err := parseThread(ioutil.NopCloser(bytes.NewBufferString(`[
			{"data": {}},
			{"data": {}}
	]`))); err == nil {
		t.Errorf("wanted error for bad link listing")
	}

	if _, err := parseThread(ioutil.NopCloser(bytes.NewBufferString(`[
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
	]`))); err == nil {
		t.Errorf("wanted error for non-one link listing")
	}

	thread, err := parseThread(ioutil.NopCloser(bytes.NewBufferString(`[
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
	]`)))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if thread.GetTitle() != "hola" {
		t.Errorf("got %s; wanted hola", thread.GetTitle())
	}
	if len(thread.GetComments()) != 2 {
		t.Errorf("got %v; wanted 2 comments", thread.GetComments)
	}

	if _, err := parseThread(ioutil.NopCloser(bytes.NewBufferString(`[
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
	]`))); err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestParseInbox(t *testing.T) {
	if _, err := parseInbox(nil); err == nil {
		t.Errorf("wanted error for nil content")
	}

	if _, err := parseInbox(
		ioutil.NopCloser(bytes.NewBufferString(`[]"`)),
	); err == nil {
		t.Errorf("wanted error for invalid json")
	}

	if _, err := parseInbox(
		ioutil.NopCloser(bytes.NewBufferString(`{}`)),
	); err == nil {
		t.Errorf("wanted error for nil data")
	}

	if _, err := parseInbox(
		ioutil.NopCloser(bytes.NewBufferString(`{"data": {}}`)),
	); err == nil {
		t.Errorf("wanted error for nil children")
	}

	messages, err := parseInbox(ioutil.NopCloser(bytes.NewBufferString(`{
		"data": {
			"children" : [
				{"data": {"was_comment": true}}
			]
		}
	}`)))
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

	topID := "top"
	subID := "sub"
	sublisting := &redditproto.CommentListing{
		Data: &redditproto.CommentData{
			Children: []*redditproto.CommentChildren{
				{
					Data: &redditproto.Comment{
						Id: &subID,
					},
				},
			},
		},
	}
	listing := &redditproto.CommentListing{
		Data: &redditproto.CommentData{
			Children: []*redditproto.CommentChildren{
				{
					Data: &redditproto.Comment{
						Id:      &topID,
						Replies: sublisting,
					},
				},
			},
		},
	}
	expected := []*redditproto.Comment{
		{
			Id: &topID,
			ReplyTree: []*redditproto.Comment{
				{
					Id: &subID,
				},
			},
		},
	}
	actual := unpackCommentListing(listing)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v; wanted %v", actual, expected)
	}
}
