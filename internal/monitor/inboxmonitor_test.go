package monitor

import (
	"fmt"
	"testing"
	"time"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

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

func TestInboxMonitorUpdate(t *testing.T) {
	im := &InboxMonitor{
		Op: &operator.MockOperator{
			InboxErr: fmt.Errorf("an error"),
		},
	}
	if err := im.Update(); err == nil {
		t.Errorf("wanted error for request failure")
	}

	mentionSubject := "username mention"
	postReplySubject := "post reply"
	tval := true
	bot := &mockInboxHandler{}
	im = &InboxMonitor{
		Op: &operator.MockOperator{
			InboxReturn: []*redditproto.Message{
				&redditproto.Message{Subject: &mentionSubject},
				&redditproto.Message{Subject: &postReplySubject},
				&redditproto.Message{WasComment: &tval},
				&redditproto.Message{},
			},
		},
		MessageHandler:      bot,
		PostReplyHandler:    bot,
		CommentReplyHandler: bot,
		MentionHandler:      bot,
	}
	if err := im.Update(); err != nil {
		t.Fatalf("error: %v", err)
	}

	// Allow bot goroutines to work.
	time.Sleep(time.Second)

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
