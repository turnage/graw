package operator

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/turnage/redditproto"
)

func TestMockOperator(t *testing.T) {
	title := "title"
	expectedErr := fmt.Errorf("an error")
	expectedScrapeReturn := []*redditproto.Link{
		&redditproto.Link{Title: &title},
	}
	expectedThreadsReturn := []*redditproto.Link{
		&redditproto.Link{Title: &title},
		&redditproto.Link{Title: &title},
	}
	expectedThreadReturn := &redditproto.Link{Title: &title}
	expectedInboxReturn := []*redditproto.Message{
		&redditproto.Message{Subject: &title},
	}

	mock := Operator(&MockOperator{
		ScrapeErr:     expectedErr,
		ScrapeReturn:  expectedScrapeReturn,
		ThreadsErr:    expectedErr,
		ThreadsReturn: expectedThreadsReturn,
		ThreadErr:     expectedErr,
		ThreadReturn:  expectedThreadReturn,
		InboxErr:      expectedErr,
		InboxReturn:   expectedInboxReturn,
		MarkAsReadErr: expectedErr,
		ReplyErr:      expectedErr,
		SubmitErr:     expectedErr,
		ComposeErr:    expectedErr,
	})

	if err := mock.MarkAsRead(); err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}

	if err := mock.Reply("", ""); err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}

	if err := mock.Compose("", "", ""); err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}

	if err := mock.Submit("", "", "", ""); err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}

	threads, err := mock.Scrape("", "", "", 0, Link)
	if !reflect.DeepEqual(threads.([]*redditproto.Link), expectedScrapeReturn) {
		t.Errorf("got %v; wanted %v", threads, expectedScrapeReturn)
	}

	if err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}

	threads, err = mock.Threads()
	if !reflect.DeepEqual(threads, expectedThreadsReturn) {
		t.Errorf("got %v; wanted %v", threads, expectedThreadsReturn)
	}

	if err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}

	thread, err := mock.Thread("")
	if !reflect.DeepEqual(thread, expectedThreadReturn) {
		t.Errorf("got %v; wanted %v", thread, expectedThreadReturn)
	}

	if err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}

	messages, err := mock.Inbox()
	if !reflect.DeepEqual(messages, expectedInboxReturn) {
		t.Errorf("got %v; wanted %v", messages, expectedInboxReturn)
	}

	if err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}
}
