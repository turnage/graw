package monitor

import (
	"reflect"
	"testing"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

type mockInboxHandler struct {
	mentionCalls      int
	postReplyCalls    int
	commentReplyCalls int
}

func (m *mockInboxHandler) Mention(msg *redditproto.Comment) {
	m.mentionCalls++
}

func (m *mockInboxHandler) PostReply(msg *redditproto.Comment) {
	m.postReplyCalls++
}

func (m *mockInboxHandler) CommentReply(msg *redditproto.Comment) {
	m.commentReplyCalls++
}

func (m *mockInboxHandler) Message(msg *redditproto.Message) {}

func TestInboxMonitor(t *testing.T) {
	mon, err := InboxMonitor(
		&operator.MockOperator{
			ScrapeMessagesReturn: []*redditproto.Message{
				&redditproto.Message{
					Name: stringPointer("name"),
				},
			},
		},
		&mockInboxHandler{},
		&mockInboxHandler{},
		&mockInboxHandler{},
		&mockInboxHandler{},
		Forward,
	)
	if err != nil {
		t.Fatal(err)
	}

	i := mon.(*inboxMonitor)

	if i.mentionHandler == nil ||
		i.postReplyHandler == nil ||
		i.commentReplyHandler == nil {
		t.Errorf("got %v; wanted all fields set", i)
	}

	if i.dir != Forward {
		t.Errorf("got %d; wanted %d (Forward)", i.dir, Forward)
	}

	if i.path != "/message/inbox" {
		t.Errorf("got %s; wanted /message/inbox", i.path)
	}

	if !reflect.DeepEqual(i.tip, []string{"name"}) {
		t.Errorf("got %v; wanted %v", i.tip, []string{"name"})
	}
}

func TestInboxMonitorCommentDispatch(t *testing.T) {
	han := &mockInboxHandler{}
	i := &inboxMonitor{
		postReplyHandler:    han,
		commentReplyHandler: han,
		mentionHandler:      han,
	}

	i.commentDispatch(
		&redditproto.Comment{
			Subject: stringPointer(mentionSubject),
		},
	)
	if han.mentionCalls != 1 {
		t.Errorf("got %d mention calls; wanted 1", han.mentionCalls)
	}

	i.commentDispatch(
		&redditproto.Comment{
			Subject: stringPointer(postReplySubject),
		},
	)
	if han.postReplyCalls != 1 {
		t.Errorf("got %d postReply calls; wanted 1", han.postReplyCalls)
	}

	i.commentDispatch(&redditproto.Comment{})
	if han.commentReplyCalls != 1 {
		t.Errorf("got %d comment calls; wanted 1", han.commentReplyCalls)
	}
}
