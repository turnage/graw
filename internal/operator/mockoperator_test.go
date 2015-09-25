package operator

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/turnage/redditproto"
)

func TestMockMarkAsRead(t *testing.T) {
	expectedErr := fmt.Errorf("an error")
	mock := Operator(
		&MockOperator{
			MarkAsReadErr: expectedErr,
		},
	)
	if err := mock.MarkAsRead(); err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}
}

func TestMockReply(t *testing.T) {
	expectedErr := fmt.Errorf("an error")
	mock := Operator(
		&MockOperator{
			ReplyErr: expectedErr,
		},
	)
	if err := mock.Reply("", ""); err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}
}

func TestMockCompose(t *testing.T) {
	expectedErr := fmt.Errorf("an error")
	mock := Operator(
		&MockOperator{
			ComposeErr: expectedErr,
		},
	)
	if err := mock.Compose("", "", ""); err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}
}

func TestMockSubmit(t *testing.T) {
	expectedErr := fmt.Errorf("an error")
	mock := Operator(
		&MockOperator{
			SubmitErr: expectedErr,
		},
	)
	if err := mock.Submit("", "", "", ""); err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}
}

func TestMockScrape(t *testing.T) {
	expectedErr := fmt.Errorf("an error")
	title := "title"
	expected := []Thing{
		&redditproto.Link{Title: &title},
	}
	mock := Operator(
		&MockOperator{
			ScrapeErr:    expectedErr,
			ScrapeReturn: expected,
		},
	)
	actual, err := mock.Scrape("", "", "", 0, Link)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v; wanted %v", actual, expected)
	}

	if err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}
}

func TestMockThreads(t *testing.T) {
	expectedErr := fmt.Errorf("an error")
	title := "title"
	expected := []*redditproto.Link{
		&redditproto.Link{Title: &title},
		&redditproto.Link{Title: &title},
	}
	mock := Operator(
		&MockOperator{
			ThreadsErr:    expectedErr,
			ThreadsReturn: expected,
		},
	)
	actual, err := mock.Threads()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v; wanted %v", actual, expected)
	}

	if err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}
}

func TestMockThread(t *testing.T) {
	expectedErr := fmt.Errorf("an error")
	title := "title"
	expected := &redditproto.Link{Title: &title}
	mock := Operator(
		&MockOperator{
			ThreadErr:    expectedErr,
			ThreadReturn: expected,
		},
	)
	actual, err := mock.Thread("")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v; wanted %v", actual, expected)
	}

	if err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}
}

func TestMockInbox(t *testing.T) {
	expectedErr := fmt.Errorf("an error")
	title := "title"
	expected := []*redditproto.Message{
		&redditproto.Message{Subject: &title},
	}
	mock := Operator(
		&MockOperator{
			InboxErr:    expectedErr,
			InboxReturn: expected,
		},
	)
	actual, err := mock.Inbox()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v; wanted %v", actual, expected)
	}

	if err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}
}
