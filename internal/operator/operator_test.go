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

func TestReplyRequest(t *testing.T) {
	if _, err := replyRequest("", "content"); err == nil {
		t.Errorf("wanted error for missing parent id")
	}

	if _, err := replyRequest("parent", ""); err == nil {
		t.Errorf("wanted error for missing content")
	}

	req, err := replyRequest("parent", "content")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if req.Method != "POST" {
		t.Errorf("got %s; wanted POST", req.Method)
	}

	if err := req.ParseForm(); err != nil {
		t.Errorf("error parsing form: %v", err)
	}

	if req.PostForm.Get("thing_id") != "parent" {
		t.Errorf("got %s; wanted parent", req.PostForm.Get("thing_id"))
	}

	if req.PostForm.Get("text") != "content" {
		t.Errorf("got %s; wanted content", req.PostForm.Get("text"))
	}
}

func TestReply(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if err := op.Reply("", ""); err == nil {
		t.Errorf("wanted error for invalid request")
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

func TestComposeRequest(t *testing.T) {
	if _, err := composeRequest("", "subject", "body"); err == nil {
		t.Errorf("wanted error for missing user")
	}

	if _, err := composeRequest("user", "", "body"); err == nil {
		t.Errorf("wanted error for missing subject")
	}

	if _, err := composeRequest("user", "subject", ""); err == nil {
		t.Errorf("wanted error for missing body")
	}

	req, err := composeRequest("user", "subject", "body")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if req.Method != "POST" {
		t.Errorf("got %s; wanted POST", req.Method)
	}

	if err := req.ParseForm(); err != nil {
		t.Errorf("error parsing form: %v", err)
	}

	if req.PostForm.Get("to") != "user" {
		t.Errorf("got %s; wanted user", req.PostForm.Get("thing_id"))
	}

	if req.PostForm.Get("subject") != "subject" {
		t.Errorf("got %s; wanted subject", req.PostForm.Get("text"))
	}

	if req.PostForm.Get("text") != "body" {
		t.Errorf("got %s; wanted body", req.PostForm.Get("text"))
	}
}

func TestCompose(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if err := op.Compose("", "", ""); err == nil {
		t.Errorf("wanted error for invalid request")
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

func TestSubmitRequest(t *testing.T) {
	if _, err := submitRequest("", "self", "title", ""); err == nil {
		t.Errorf("wanted error for missing subreddit")
	}

	if _, err := submitRequest("aww", "wrong", "title", ""); err == nil {
		t.Errorf("wanted error for unsupported post type")
	}

	if _, err := submitRequest("aww", "link", "", "url"); err == nil {
		t.Errorf("wanted error for omitted title")
	}

	if _, err := submitRequest("aww", "link", "title", ""); err == nil {
		t.Errorf("wanted error for omitted url")
	}

	req, err := submitRequest("aww", "self", "title", "mombo")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if req.Method != "POST" {
		t.Errorf("got %s; wanted POST", req.Method)
	}

	if err := req.ParseForm(); err != nil {
		t.Errorf("error parsing form: %v", err)
	}

	if req.PostForm.Get("sr") != "aww" {
		t.Errorf("got %s; wanted saww", req.PostForm.Get("sr"))
	}

	if req.PostForm.Get("kind") != "self" {
		t.Errorf("got %s; wanted self", req.PostForm.Get("kind"))
	}

	if req.PostForm.Get("title") != "title" {
		t.Errorf("got %s; wanted title", req.PostForm.Get("title"))
	}

	if req.PostForm.Get("text") != "mombo" {
		t.Errorf("got %s; wanted mombo", req.PostForm.Get("text"))
	}

	req, err = submitRequest("aww", "link", "title", "mombo")
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if err := req.ParseForm(); err != nil {
		t.Errorf("error parsing form: %v", err)
	}

	if req.PostForm.Get("url") != "mombo" {
		t.Errorf("got %s; wanted mombo", req.PostForm.Get("url"))
	}

	if req.PostForm.Get("kind") != "link" {
		t.Errorf("got %s; wanted link", req.PostForm.Get("kind"))
	}
}

func TestSubmit(t *testing.T) {
	op := &operator{
		cli: client.NewMock("", fmt.Errorf("an error")),
	}

	if err := op.Submit("", "", "", ""); err == nil {
		t.Errorf("wanted error for invalid request")
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
