package monitor

import (
	"fmt"
	"testing"
	"time"

	"github.com/turnage/graw/internal/monitor/internal/scanner"
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
		&operator.MockOperator{},
		&mockPostHandler{},
		[]string{"self"},
	)
	mon := pm.(*postMonitor)
	mon.postScanner = &scanner.MockScanner{
		ScanErr: fmt.Errorf("an error"),
	}
	if err := pm.Update(); err == nil {
		t.Errorf("wanted error for request failure")
	}

	bot := &mockPostHandler{}
	pm = PostMonitor(
		&operator.MockOperator{},
		bot,
		[]string{"self"},
	)
	mon = pm.(*postMonitor)
	mon.postScanner = &scanner.MockScanner{
		ScanLinksReturn: []*redditproto.Link{
			&redditproto.Link{},
			&redditproto.Link{},
		},
	}
	if err := pm.Update(); err != nil {
		t.Fatalf("error: %v", err)
	}

	for i := 0; i < 100 && bot.Calls < 2; i++ {
		time.Sleep(10 * time.Millisecond)
	}

	if bot.Calls != 2 {
		t.Errorf("%d calls were made to mock bot; wanted 1", bot.Calls)
	}
}
