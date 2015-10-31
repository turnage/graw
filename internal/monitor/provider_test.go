package monitor

import (
	"testing"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

var (
	mockForSync = &operator.MockOperator{
		ScrapeLinksReturn: []*redditproto.Link{
			&redditproto.Link{
				Name: stringPointer("name"),
			},
		},
	}
)

type mockPostHandler struct{}

func (m *mockPostHandler) Post(post *redditproto.Link) {}

type mockUserHandler struct{}

func (m *mockUserHandler) UserPost(post *redditproto.Link) {}

func (m *mockUserHandler) UserComment(comment *redditproto.Comment) {}

type mockInboxHandler struct{}

func (m *mockInboxHandler) Mention(msg *redditproto.Comment) {}

func (m *mockInboxHandler) PostReply(msg *redditproto.Comment) {}

func (m *mockInboxHandler) CommentReply(msg *redditproto.Comment) {}

func (m *mockInboxHandler) Message(msg *redditproto.Message) {}

func TestPostMonitor(t *testing.T) {
	mon, err := PostMonitor(
		mockForSync,
		&mockPostHandler{},
		[]string{"self", "aww"},
		Forward,
	)
	if err != nil {
		t.Fatal(err)
	}

	pm := mon.(*base)
	if pm.handlePost == nil {
		t.Errorf("wanted post handler set")
	}
	if pm.path != "/r/self+aww" {
		t.Errorf("got %s; wanted /r/self+aww", pm.path)
	}
}

func TestUserMonitor(t *testing.T) {
	mon, err := UserMonitor(
		mockForSync,
		&mockUserHandler{},
		"rob",
		Forward,
	)
	if err != nil {
		t.Fatal(err)
	}

	u := mon.(*base)
	if u.handlePost == nil {
		t.Errorf("wanted post handler set")
	}
	if u.handleComment == nil {
		t.Errorf("wanted comment handler set")
	}
	if u.path != "/user/rob" {
		t.Errorf("got %s; wanted /user/rob", u.path)
	}
}

func TestMessageMonitor(t *testing.T) {
	mon, err := MessageMonitor(
		mockForSync,
		&mockInboxHandler{},
		Forward,
	)
	if err != nil {
		t.Fatal(err)
	}

	m := mon.(*base)
	if m.handleMessage == nil {
		t.Errorf("wanted message handler set")
	}
	if m.path != "/message/messages" {
		t.Errorf("got %s; wanted /message/messages", m.path)
	}
}

func TestCommentReplyMonitor(t *testing.T) {
	mon, err := CommentReplyMonitor(
		mockForSync,
		&mockInboxHandler{},
		Forward,
	)
	if err != nil {
		t.Fatal(err)
	}

	c := mon.(*base)
	if c.handleComment == nil {
		t.Errorf("wanted comment handler set")
	}
	if c.path != "/message/comments" {
		t.Errorf("got %s; wanted /message/comments", c.path)
	}
}

func TestPostReplyMonitor(t *testing.T) {
	mon, err := PostReplyMonitor(
		mockForSync,
		&mockInboxHandler{},
		Forward,
	)
	if err != nil {
		t.Fatal(err)
	}

	p := mon.(*base)
	if p.handleComment == nil {
		t.Errorf("wanted comment handler set")
	}
	if p.path != "/message/selfreply" {
		t.Errorf("got %s; wanted /message/selfreply", p.path)
	}
}

func TestMentionMonitor(t *testing.T) {
	mon, err := MentionMonitor(
		mockForSync,
		&mockInboxHandler{},
		Forward,
	)
	if err != nil {
		t.Fatal(err)
	}

	m := mon.(*base)
	if m.handleComment == nil {
		t.Errorf("wanted comment handler set")
	}
	if m.path != "/message/mentions" {
		t.Errorf("got %s; wanted /message/mentions", m.path)
	}
}
