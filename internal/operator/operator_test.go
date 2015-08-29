package operator

import (
	"fmt"
	"testing"

	"github.com/turnage/graw/internal/operator/internal/client"
)

func TestScrape(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if _, err := op.Scrape("self", "new", "", "", 1); err == nil {
		t.Errorf("wanted error for request error")
	}

	op = &operator{
		cli: client.NewMock(
			`{
				"data": {
					"children": [
						{"data": {"title": "hello"}},
						{"data": {"title": "hola"}},
						{"data": {"title": "bye"}}
					]
				}
			}`,
			nil,
		),
	}

	posts, err := op.Scrape("self", "new", "", "", 3)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(posts) != 3 {
		t.Errorf("got %d posts; wanted 3", len(posts))
	}
}

func TestThreads(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if _, err := op.Threads("1", "2", "3"); err == nil {
		t.Errorf("wanted error for request failure")
	}

	op = &operator{
		cli: client.NewMock(
			`{
				"data": {
					"children": [
						{"data": {"title": "hello"}},
						{"data": {"title": "hola"}},
						{"data": {"title": "bye"}}
					]
				}
			}`,
			nil,
		),
	}

	posts, err := op.Threads("1", "2", "3")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(posts) != 3 {
		t.Errorf("got %d posts; wanted 3", len(posts))
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
