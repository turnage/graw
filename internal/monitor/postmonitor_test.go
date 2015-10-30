package monitor

import (
	"testing"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

type mockPostHandler struct {
	Calls int
}

func (m *mockPostHandler) Post(post *redditproto.Link) {
	m.Calls++
}

func TestPostMonitor(t *testing.T) {
	mon, err := PostMonitor(
		&operator.MockOperator{
			ScrapeLinksReturn: []*redditproto.Link{
				&redditproto.Link{
					Name: stringPointer("name"),
				},
			},
		},
		&mockPostHandler{},
		[]string{"self", "aww"},
		Forward,
	)
	if err != nil {
		t.Fatal(err)
	}

	pm := mon.(*postMonitor)
	if pm.dir != Forward {
		t.Errorf("got %d; wanted %d (Forward)", pm.dir, Forward)
	}
	if pm.handlePost == nil {
		t.Errorf("wanted post handler set")
	}
	if pm.path != "/r/self+aww" {
		t.Errorf("got %s; wanted /r/self+aww", pm.path)
	}
}
