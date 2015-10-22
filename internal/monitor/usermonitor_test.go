package monitor

import (
	"fmt"
	"testing"
	"time"

	"github.com/turnage/graw/internal/monitor/internal/scanner"
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

type mockUserHandler struct {
	postCalls    int
	commentCalls int
}

func (m *mockUserHandler) UserPost(post *redditproto.Link) {
	m.postCalls++
}

func (m *mockUserHandler) UserComment(comment *redditproto.Comment) {
	m.commentCalls++
}

func TestUserMonitor(t *testing.T) {
	if um := UserMonitor(
		&operator.MockOperator{},
		&mockNoHandler{},
		"user",
	); um != nil {
		t.Errorf("got %v; wanted nil", um)
	}

	if um := UserMonitor(
		&operator.MockOperator{},
		&mockUserHandler{},
		"",
	); um != nil {
		t.Errorf("got %v; wanted nil", um)
	}

	um := UserMonitor(
		&operator.MockOperator{},
		&mockUserHandler{},
		"user",
	).(*userMonitor)
	if um.userHandler == nil {
		t.Errorf("wanted userHandler set")
	}
}

func TestUserMonitorUpdate(t *testing.T) {
	um := UserMonitor(
		&operator.MockOperator{},
		&mockUserHandler{},
		"user",
	)
	mon := um.(*userMonitor)
	mon.userScanner = &scanner.MockScanner{
		ScanErr: fmt.Errorf("an error"),
	}
	if err := um.Update(); err == nil {
		t.Errorf("wanted error for request failure")
	}

	bot := &mockUserHandler{}
	um = UserMonitor(
		&operator.MockOperator{},
		bot,
		"user",
	)
	mon = um.(*userMonitor)
	mon.userScanner = &scanner.MockScanner{
		ScanLinksReturn: []*redditproto.Link{
			&redditproto.Link{},
			&redditproto.Link{},
		},
		ScanCommentsReturn: []*redditproto.Comment{
			&redditproto.Comment{},
			&redditproto.Comment{},
		},
	}
	if err := um.Update(); err != nil {
		t.Fatalf("error: %v", err)
	}

	for i := 0; i < 100 && bot.postCalls+bot.commentCalls < 4; i++ {
		time.Sleep(10 * time.Millisecond)
	}

	if bot.postCalls != 2 {
		t.Errorf(
			"%d calls were made to mock bot; wanted 2",
			bot.postCalls,
		)
	}

	if bot.commentCalls != 2 {
		t.Errorf(
			"%d calls were made to mock bot; wanted 2",
			bot.commentCalls,
		)
	}
}
