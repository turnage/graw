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
		t.Errorf("got %v; wanted nil")
	}

	if pm := PostMonitor(
		&operator.MockOperator{},
		&mockPostHandler{},
		[]string{},
	); pm != nil {
		t.Errorf("got %v; wanted nil")
	}

	pm := PostMonitor(
		&operator.MockOperator{},
		&mockPostHandler{},
		[]string{"self", "aww"},
	).(*postMonitor)
	if pm.postHandler == nil {
		t.Errorf("wanted postHandler set")
	}
	if pm.query != "self+aww" {
		t.Errorf("got %s; wanted self+aww", pm.query)
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

	pm = PostMonitor(
		&operator.MockOperator{
			ThreadsErr: fmt.Errorf("an error"),
		},
		&mockPostHandler{},
		[]string{"self"},
	)
	if err := pm.Update(); err == nil {
		t.Errorf("wanted error for fixtip failure")
	}

	bot := &mockPostHandler{}
	postName := "name"
	pm = PostMonitor(
		&operator.MockOperator{
			ScrapeReturn: []*redditproto.Link{
				&redditproto.Link{Name: &postName},
				&redditproto.Link{Name: &postName},
			},
			ThreadsReturn: []*redditproto.Link{
				&redditproto.Link{Name: &postName},
			},
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

func TestFetchTip(t *testing.T) {
	pm := &postMonitor{
		tip: []string{""},
	}

	pm.op = &operator.MockOperator{
		ScrapeErr: fmt.Errorf("an error"),
	}
	if _, err := pm.fetchTip(); err == nil {
		t.Errorf("wanted error for request failure")
	}

	pm.tip = make([]string, maxTipSize)
	for i := 0; i < maxTipSize; i++ {
		pm.tip = append(pm.tip, "id")
	}
	postName := "anything"
	pm.op = &operator.MockOperator{
		ScrapeErr: nil,
		ScrapeReturn: []*redditproto.Link{
			&redditproto.Link{Name: &postName},
		},
	}

	posts, err := pm.fetchTip()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if pm.tip[len(pm.tip)-1] != postName {
		t.Errorf(
			"got tip %s; wanted %s",
			pm.tip[len(pm.tip)-1],
			postName)
	}

	if len(posts) != 1 {
		t.Fatalf("got %d posts; expected 1", len(posts))
	}

	if posts[0].GetName() != postName {
		t.Errorf(
			"got thread name %s; wanted %s",
			posts[0].GetName(),
			postName)
	}

	pm.tip = []string{""}
	posts, err = pm.fetchTip()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if posts != nil {
		t.Errorf("got %v; wanted no posts for adjustment round", posts)
	}
}

func TestFixTip(t *testing.T) {
	pm := &postMonitor{
		tip: []string{"1", "2", "3"},
	}

	pm.op = &operator.MockOperator{
		ThreadsErr: fmt.Errorf("an error"),
	}
	if err := pm.fixTip(); err == nil {
		t.Errorf("wanted error for request failure")
	}

	pm.op = &operator.MockOperator{
		ThreadsErr:    nil,
		ThreadsReturn: nil,
	}
	if err := pm.fixTip(); err != nil {
		t.Fatalf("error: %v", err)
	}

	if pm.tip[len(pm.tip)-1] != "2" {
		t.Errorf(
			"got %s; wanted tip shaved to 2",
			pm.tip[len(pm.tip)-1])
	}
}

func TestShaveTip(t *testing.T) {
	pm := &postMonitor{
		tip: []string{"1", "2"},
	}

	pm.shaveTip()
	if pm.tip[len(pm.tip)-1] != "1" {
		t.Errorf(
			"got %s; wanted tip shaved to 1",
			pm.tip[len(pm.tip)-1])
	}

	pm.shaveTip()
	if len(pm.tip) != 1 {
		t.Errorf("tip is %d long; wanted 1 blank tip", len(pm.tip))
	}

	if pm.tip[0] != "" {
		t.Errorf("got %s; wanted empty string", pm.tip[0])
	}
}
