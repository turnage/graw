package monitor

import (
	"fmt"
	"testing"
	"time"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

type mockNoHandler struct{}

type mockInboxHandler struct {
	MentionCalls      int
	PostReplyCalls    int
	CommentReplyCalls int
	MessageCalls      int
}

func (m *mockInboxHandler) Mention(msg *redditproto.Message) {
	m.MentionCalls++
}

func (m *mockInboxHandler) PostReply(msg *redditproto.Message) {
	m.PostReplyCalls++
}

func (m *mockInboxHandler) CommentReply(msg *redditproto.Message) {
	m.CommentReplyCalls++
}

func (m *mockInboxHandler) Message(msg *redditproto.Message) {
	m.MessageCalls++
}

func TestInboxMonitor(t *testing.T) {
	expectedOperator := &operator.MockOperator{}
	if im := InboxMonitor(
		expectedOperator,
		&mockNoHandler{},
	); im != nil {
		t.Errorf("got %v; wanted nil", im)
	}

	im := InboxMonitor(
		expectedOperator,
		&mockInboxHandler{},
	).(*inboxMonitor)

	if im.op != expectedOperator {
		t.Errorf("got %v; wanted %v", im.op, expectedOperator)
	}

	if im.messageHandler == nil ||
		im.mentionHandler == nil ||
		im.postReplyHandler == nil ||
		im.commentReplyHandler == nil {
		t.Errorf("got %v; wanted all fields set")
	}
}

func TestInboxMonitorUpdate(t *testing.T) {
	im := InboxMonitor(
		&operator.MockOperator{
			InboxErr: fmt.Errorf("an error"),
		},
		&mockInboxHandler{},
	)
	if err := im.Update(); err == nil {
		t.Errorf("wanted error for request failure")
	}

	mentionSubject := "username mention"
	postReplySubject := "post reply"
	tval := true
	bot := &mockInboxHandler{}
	im = InboxMonitor(
		&operator.MockOperator{
			InboxReturn: []*redditproto.Message{
				&redditproto.Message{Subject: &mentionSubject},
				&redditproto.Message{Subject: &postReplySubject},
				&redditproto.Message{WasComment: &tval},
				&redditproto.Message{},
			},
		},
		bot,
	)
	if err := im.Update(); err != nil {
		t.Fatalf("error: %v", err)
	}

	// Allow bot goroutines to work.
	time.Sleep(20 * time.Millisecond)

	if bot.MentionCalls != 1 {
		t.Errorf("wanted a call to Mention()")
	}

	if bot.PostReplyCalls != 1 {
		t.Errorf("wanted a call to PostReply()")
	}

	if bot.CommentReplyCalls != 1 {
		t.Errorf("wanted a call to CommentReply()")
	}

	if bot.MessageCalls != 1 {
		t.Errorf("wanted a call to Message()")
	}
}
