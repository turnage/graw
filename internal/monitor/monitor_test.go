package monitor

import (
	"container/list"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

var (
	names        = []string{"1", "2", "3", "4"}
	mockOperator = &operator.MockOperator{
		ScrapeReturn: []*redditproto.Link{
			&redditproto.Link{Name: &names[3]},
			&redditproto.Link{Name: &names[2]},
			&redditproto.Link{Name: &names[1]},
			&redditproto.Link{Name: &names[0]},
		},
		ThreadsReturn: []*redditproto.Link{
			&redditproto.Link{Name: &names[0]},
			&redditproto.Link{Name: &names[1]},
			&redditproto.Link{Name: &names[2]},
			&redditproto.Link{Name: &names[3]},
		},
	}
)

func TestNew(t *testing.T) {
	if mon := New(
		&operator.MockOperator{},
		[]string{"test"},
	); mon.NewPosts == nil ||
		mon.Errors == nil ||
		mon.errorBackOffUnit == 0 ||
		mon.op == nil ||
		mon.tip == nil ||
		mon.monitoredSubreddits == nil ||
		mon.subredditQuery == "" {
		t.Errorf("there was an uninitialized field: %v", mon)
	}
}

func TestMonitorToggles(t *testing.T) {
	mon := &Monitor{
		monitoredSubreddits: make(map[string]bool),
	}

	mon.MonitorSubreddits("awww", "self")
	if !strings.Contains(
		mon.subredditQuery,
		"aww",
	) || !strings.Contains(
		mon.subredditQuery,
		"self",
	) {
		t.Errorf(
			"got %s; wanted awww+self (any order of)",
			mon.subredditQuery)
	}
	mon.UnmonitorSubreddits("awww")
	if mon.subredditQuery != "self" {
		t.Errorf("got %s; wanted self", mon.subredditQuery)
	}
}

func TestCheckOnTip(t *testing.T) {
	mon := &Monitor{
		op:  mockOperator,
		tip: list.New(),
	}
	mon.tip.PushFront("4")

	if err := mon.checkOnTip(1); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	if mon.blanks != 0 {
		t.Errorf("unwanted blank count increment; got %d", mon.blanks)
	}

	oldTolerance := mon.blankRoundTolerance
	mon.blanks = mon.blankRoundTolerance
	if err := mon.checkOnTip(0); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	if mon.blanks != 0 {
		t.Errorf("got blank round count %d; wanted 0", mon.blanks)
	}
	if mon.blankRoundTolerance <= oldTolerance {
		t.Errorf(
			"got %d tolerance; wanted > %d",
			mon.blankRoundTolerance,
			oldTolerance)
	}
}

func TestErrorBackoff(t *testing.T) {
	mon := &Monitor{
		Errors: make(chan error),
	}

	mon.errors = errorTolerance - 1
	mon.errorBackOff(nil)
	if mon.errors != 0 {
		t.Errorf("got error count %d; wanted it zero'd", mon.errors)
	}

	if mon.errorBackOff(fmt.Errorf("an error")) {
		t.Errorf("should have forgiven error within tolerance")
	}

	mon.errors = errorTolerance
	go func() { <-mon.Errors }()
	if !mon.errorBackOff(fmt.Errorf("an error")) {
		t.Errorf("wanted indication that error tolerance was exceeded")
	}
}

func TestUpdatePosts(t *testing.T) {
	mon := &Monitor{
		op:       mockOperator,
		NewPosts: make(chan *redditproto.Link),
		tip:      list.New(),
	}
	mon.tip.PushFront("")

	go func() {
		for true {
			<-mon.NewPosts
		}
	}()

	count, err := mon.updatePosts()
	if err != nil {
		t.Fatal("error: %v", err)
	}

	if count != 4 {
		t.Errorf("got %d posts, wanted 4 posts", count)
	}
}

func TestTip(t *testing.T) {
	mon := &Monitor{
		op:  mockOperator,
		tip: list.New(),
	}
	mon.tip.PushFront("shouldnotbeatback")
	for i := 0; i < maxTipSize-1; i++ {
		mon.tip.PushFront("bunk")
	}

	posts, err := mon.fetchTip()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if mon.tip.Back().Value.(string) == "shouldnotbeatback" {
		t.Errorf("tips were not truncated at capacity")
	}

	if mon.tip.Front().Value.(string) != "4" {
		t.Errorf(
			"wanted front tip '4'; got %s",
			mon.tip.Front().Value.(string))
	}

	if len(posts) != 4 {
		t.Errorf("wanted 4 incrementally named posts; got %v", posts)
	}
}

func TestFixTip(t *testing.T) {
	mon := &Monitor{
		op:  mockOperator,
		tip: list.New(),
	}
	mon.tip.PushFront("1")
	mon.tip.PushFront("0")
	broken, err := mon.fixTip()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if mon.tip.Front().Value.(string) != "1" {
		t.Errorf(
			"wanted '1'; got %s",
			mon.tip.Front().Value.(string))
	}
	if !broken {
		t.Errorf("wanted fixTip to indicate broken; 0 was gone")
	}

	mon.tip = list.New()
	mon.tip.PushFront("mark")
	mon.tip.PushFront("internet")
	_, err = mon.fixTip()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if mon.tip.Front().Value.(string) != "" {
		t.Errorf(
			"wanted ''; got %s",
			mon.tip.Front().Value.(string))
	}

	mon.tip = list.New()
	mon.tip.PushFront("1")
	broken, err = mon.fixTip()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if broken {
		t.Errorf("wanted broken to be false; '1' was valid tip")
	}
}

func TestSetKeys(t *testing.T) {
	expected := map[string]bool{
		"1": false,
		"2": false,
		"3": true,
		"4": false,
	}
	switches := map[string]bool{
		"1": true,
		"2": true,
		"3": true,
	}
	setKeys(switches, false, []string{"1", "2", "4"})
	if !reflect.DeepEqual(switches, expected) {
		t.Errorf("got %v; wanted %s", switches, expected)
	}
}

func TestBuildQuery(t *testing.T) {
	if actual := buildQuery(
		map[string]bool{
			"rob":   true,
			"joe":   true,
			"nikko": false,
		},
		"-",
	); actual != "rob-joe" && actual != "joe-rob" {
		t.Errorf("got %s; wanted rob-joe or joe-rob", actual)
	}

	if actual := buildQuery(nil, "-"); actual != "" {
		t.Errorf("got %s; wanted empty string", actual)
	}
}
