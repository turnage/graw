package monitor

import (
	"testing"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

type mockUserHandler struct{}

func (m *mockUserHandler) UserPost(post *redditproto.Link) {}

func (m *mockUserHandler) UserComment(comment *redditproto.Comment) {}

func TestUserMonitor(t *testing.T) {
	mon, err := UserMonitor(
		&operator.MockOperator{
			ScrapeLinksReturn: []*redditproto.Link{
				&redditproto.Link{},
			},
		},
		&mockUserHandler{},
		"rob",
		Forward,
	)
	if err != nil {
		t.Fatal(err)
	}

	u := mon.(*userMonitor)
	if u.handlePost == nil {
		t.Errorf("wanted post handler set")
	}
	if u.handleComment == nil {
		t.Errorf("wanted comment handler set")
	}
	if u.path != "/user/rob" {
		t.Errorf("got %s; wanted /user/rob", u.path)
	}
	if u.dir != Forward {
		t.Errorf("got %d; wanted %d (Forward)", u.dir, Forward)
	}
}
