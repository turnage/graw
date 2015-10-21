package operator

import (
	"fmt"
	"testing"

	"github.com/turnage/graw/internal/operator/internal/client"
)

var (
	linkListingJSON = `{
		"kind": "Listing",
		"data": {
			"children": [
				{
					"kind": "t3",
					"data": {
						"title": "hello",
						"body": "hello"
					}
				}
			]
		}
	}`
	comboListingJSON = `{
		"kind": "Listing",
		"data": {
			"children": [
				{
					"kind": "t3",
					"data": {
						"title": "hello",
						"body": "hello"
					}
				},
				{
					"kind": "t1",
					"data": {
						"body": "hello"
					}
				}
			]
		}
	}`
	threadJSON = `[
		{
			"kind": "Listing",
			"data": {
				"children": [
					{
						"kind": "t3",
						"data": {
							"title": "hola"
						}
					}
				]
			}
		},
		{
			"kind": "Listing",
			"data": {
				"children": [
					{
						"kind": "t1",
						"data": {
							"id": "arnold"
						}
					},
					{
						"kind": "t1",
						"data": {
							"id": "harold"
						}
					}
				]
			}
		}
	]`
	inboxJSON = `{
		"kind": "Listing",
		"data": {
			"children" : [
				{
					"kind": "t4",
					"data": {
						"was_comment": true
					}
				}
			]
		}
	}`
)

func TestPosts(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if _, err := op.Posts("/r/self/new", "", "", 1); err == nil {
		t.Errorf("wanted error for request error")
	}

	op = &operator{
		cli: client.NewMock(linkListingJSON, nil),
	}

	posts, err := op.Posts("self", "", "", 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 1 {
		t.Errorf("got %d posts; wanted 1", len(posts))
	}
}

func TestUserContent(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if _, _, err := op.UserContent("user", "", "", 1); err == nil {
		t.Errorf("wanted error for request error")
	}

	op = &operator{
		cli: client.NewMock(comboListingJSON, nil),
	}

	posts, comments, err := op.UserContent("user", "", "", 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 1 {
		t.Errorf("got %d posts; wanted 1", len(posts))
	}

	if len(comments) != 1 {
		t.Errorf("got %d comments; wanted 1", len(comments))
	}
}

func TestIsThereThing(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if _, err := op.IsThereThing("1"); err == nil {
		t.Errorf("wanted error for request failure")
	}

	op = &operator{
		cli: client.NewMock(linkListingJSON, nil),
	}

	exists, err := op.IsThereThing("1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if exists == false {
		t.Errorf("got false; wanted true")
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
	_, err = op.IsThereThing("1")
	if err == nil {
		t.Errorf("wanted an error for a bad response")
	}

	// Missing Things should return nil.
	op = &operator{
		cli: client.NewMock(
			`{
				"kind": "Listing",
				"data": {
					"children": []
				}
			}`,
			nil,
		),
	}
	exists, err = op.IsThereThing("1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if exists != false {
		t.Errorf("got true; wanted false")
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
		cli: client.NewMock(threadJSON, nil),
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
		cli: client.NewMock(inboxJSON, nil),
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
