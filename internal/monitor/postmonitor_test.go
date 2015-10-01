package monitor

import (
	"fmt"
	"testing"
	"time"

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
	if pm := PostMonitor(
		&operator.MockOperator{},
		&mockNoHandler{},
		[]string{"self"},
	); pm != nil {
		t.Errorf("got %v; wanted nil", pm)
	}

	if pm := PostMonitor(
		&operator.MockOperator{},
		&mockPostHandler{},
		[]string{},
	); pm != nil {
		t.Errorf("got %v; wanted nil", pm)
	}

	pm := PostMonitor(
		&operator.MockOperator{},
		&mockPostHandler{},
		[]string{"self", "aww"},
	).(*postMonitor)
	if pm.postHandler == nil {
		t.Errorf("wanted postHandler set")
	}
}

func TestPostMonitorUpdate(t *testing.T) {
	pm := PostMonitor(
		&operator.MockOperator{
			ScrapeErr: fmt.Errorf("an error"),
		},
		&mockPostHandler{},
		[]string{"self"},
	)
	if err := pm.Update(); err == nil {
		t.Errorf("wanted error for request failure")
	}

	bot := &mockPostHandler{}
	postName := "name"
	pm = PostMonitor(
		&operator.MockOperator{
			ScrapeReturn: []operator.Thing{
				&redditproto.Link{Name: &postName},
				&redditproto.Link{Name: &postName},
			},
			GetThingReturn: &redditproto.Link{Name: &postName},
		},
		bot,
		[]string{"self"},
	)
	if err := pm.Update(); err != nil {
		t.Fatalf("error: %v", err)
	}

	// Allow bot goroutines to work.
	time.Sleep(time.Second)

	if bot.Calls != 2 {
		t.Errorf("%d calls were made to mock bot; wanted 1", bot.Calls)
	}
}
