package monitor

import (
	"testing"

	"github.com/turnage/redditproto"
)

var (
	mockForSync = MockScraper(
		[]*redditproto.Link{
			&redditproto.Link{
				Name: stringPointer("name"),
			},
		}, nil, nil, nil,
	)
)

func post(post *redditproto.Link)          {}
func comment(comment *redditproto.Comment) {}
func message(message *redditproto.Message) {}

func TestPostMonitor(t *testing.T) {
	mon, err := PostMonitor(
		mockForSync,
		post,
		[]string{"self", "aww"},
	)
	if err != nil {
		t.Fatal(err)
	}

	pm := mon.(*base)
	if pm.handlePost == nil {
		t.Errorf("wanted post handler set")
	}
	if pm.path != "/r/self+aww/new" {
		t.Errorf("got %s; wanted /r/self+aww/new", pm.path)
	}
}

func TestUserMonitor(t *testing.T) {
	mon, err := UserMonitor(
		mockForSync,
		post,
		comment,
		"rob",
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
		message,
	)
	if err != nil {
		t.Fatal(err)
	}

	m := mon.(*messageMonitor)
	if m.handleMessage == nil {
		t.Errorf("wanted message handler set")
	}
	if m.path != "/message/inbox" {
		t.Errorf("got %s; wanted /message/inbox", m.path)
	}
}

func TestCommentReplyMonitor(t *testing.T) {
	mon, err := CommentReplyMonitor(
		mockForSync,
		comment,
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
		comment,
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
		comment,
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
