package operator

import (
	"fmt"
	"testing"

	"github.com/turnage/graw/internal/operator/internal/client"
	"github.com/turnage/redditproto"
)

func TestScrape(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if _, err := op.Scrape("/r/self/new", "", "", 1, Link); err == nil {
		t.Errorf("wanted error for request error")
	}

	op = &operator{
		cli: client.NewMock(
			`{
				"data": {
					"children": [
						{"data": {
							"title": "hello",
							"body": "hello"
						}},
						{"data": {
							"title": "hola",
							"body": "hola"
						}},
						{"data": {
							"title": "bye",
							"body": "bye"
						}}
					]
				}
			}`,
			nil,
		),
	}

	postThings, err := op.Scrape("/r/self/new", "", "", 1, Link)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	posts := make([]*redditproto.Link, len(postThings))
	for i, thing := range postThings {
		posts[i] = thing.(*redditproto.Link)
	}

	if len(posts) != 3 {
		t.Errorf("got %d posts; wanted 3", len(posts))
	}

	if posts[0].GetTitle() != "hello" {
		t.Errorf("got %s; wanted hello", posts[0].GetTitle())
	}

	commentThings, err := op.Scrape("/r/self/new", "", "", 1, Comment)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	comments := make([]*redditproto.Comment, len(commentThings))
	for i, thing := range commentThings {
		comments[i] = thing.(*redditproto.Comment)
	}

	if comments[0].GetBody() != "hello" {
		t.Errorf("got %s; wanted hello", comments[0].GetBody())
	}
}

func TestGetThing(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if _, err := op.GetThing("1", Link); err == nil {
		t.Errorf("wanted error for request failure")
	}

	op = &operator{
		cli: client.NewMock(
			`{
				"data": {
					"children": [
						{"data": {"name": "charlie"}}
					]
				}
			}`,
			nil,
		),
	}

	thing, err := op.GetThing("1", Link)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if thing == nil {
		t.Fatalf("wanted a thing returned")
	}

	if thing.GetName() != "charlie" {
		t.Errorf("got %s posts; wanted charlie", thing.GetName())
	}

	// Bad responses should bubble an error up.
	op = &operator{
		cli: client.NewMock(
			`{
				"data": {
					"children": [
						{"data": {"name": "charlie"}},
					]
				}
			}`,
			nil,
		),
	}
	thing, err = op.GetThing("1", Link)
	if err == nil {
		t.Errorf("wanted an error for a bad response")
	}

	// Missing Things should return nil.
	op = &operator{
		cli: client.NewMock(
			`{
				"data": {
					"children": []
				}
			}`,
			nil,
		),
	}
	thing, err = op.GetThing("1", Link)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if thing != nil {
		t.Errorf("got %v; wanted nil", thing)
	}
}

func TestThread(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if _, err := op.Thread("/thread"); err == nil {
		t.Errorf("wanted error for request failure")
	}

	op = &operator{
		cli: client.NewMock(`[
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
		]`, nil),
	}

	thread, err := op.Thread("/thread")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(thread.GetComments()) != 2 {
		t.Errorf("got %d comments; wanted 2", len(thread.GetComments()))
	}
}

func TestInbox(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if _, err := op.Inbox(); err == nil {
		t.Errorf("wanted error for request failure")
	}

	op = &operator{
		cli: client.NewMock(`{
			"data": {
				"children" : [
					{"data": {"was_comment": true}}
				]
			}
		}`, nil),
	}

	messages, err := op.Inbox()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("got %d messages; wanted 1", len(messages))
	}

	if !messages[0].GetWasComment() {
		t.Fatal("got non-comment inboxable; wanted comment inboxable")
	}
}

func TestMarkAsRead(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if err := op.MarkAsRead("id1", "id2"); err == nil {
		t.Errorf("wanted error for request failure")
	}

	op = &operator{
		cli: client.NewMock("", nil),
	}

	if err := op.MarkAsRead("id1", "id2"); err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestReply(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if err := op.Reply("parent", "content"); err == nil {
		t.Errorf("wanted error for request failure")
	}

	op = &operator{
		cli: client.NewMock("", nil),
	}

	if err := op.Reply("parent", "content"); err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestCompose(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if err := op.Compose("user", "subject", "body"); err == nil {
		t.Errorf("wanted error for request failure")
	}

	op = &operator{
		cli: client.NewMock("", nil),
	}

	if err := op.Compose("user", "subject", "body"); err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestSubmit(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if err := op.Submit("aww", "self", "title", ""); err == nil {
		t.Errorf("wanted error for request failure")
	}

	op = &operator{
		cli: client.NewMock("", nil),
	}

	if err := op.Submit("aww", "self", "title", ""); err != nil {
		t.Fatalf("error: %v", err)
	}
}
