package monitor

import (
	"container/list"
	"reflect"
	"strings"
	"testing"

	"github.com/turnage/graw/bot/internal/client"
	"github.com/turnage/graw/bot/internal/operator"
)

func TestMonitorToggles(t *testing.T) {
	mon := &Monitor{
		monitoredSubreddits: make(map[string]bool),
		monitoredThreads:    make(map[string]bool),
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

	mon.MonitorThreads("harry", "potter")
	if !strings.Contains(
		mon.threadQuery,
		"harry",
	) || !strings.Contains(
		mon.threadQuery,
		"potter",
	) {
		t.Errorf(
			"got %s; wanted harry,potter (any order of)",
			mon.threadQuery)
	}
	mon.UnmonitorThreads("potter")
	if mon.threadQuery != "harry" {
		t.Errorf("got %s; wanted harry", mon.threadQuery)
	}
}

func TestTip(t *testing.T) {
	mon := &Monitor{
		op: operator.New(
			client.NewMock(`{
				"data": {
					"children": [
						{"data":{"name":"4"}},
						{"data":{"name":"3"}},
						{"data":{"name":"2"}},
						{"data":{"name":"1"}}
					]
				}
			}`),
		),
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
		op: operator.New(
			client.NewMock(`{
				"data": {
					"children": [
						{"data":{"name":"1"}},
						{"data":{"name":"2"}},
						{"data":{"name":"3"}},
						{"data":{"name":"4"}}
					]
				}
			}`),
		),
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
	expected := "rob-joe"
	if actual := buildQuery(
		map[string]bool{
			"rob":   true,
			"joe":   true,
			"nikko": false,
		},
		"-",
	); actual != expected {
		t.Errorf("got %s; wanted %s", actual, expected)
	}

	if actual := buildQuery(nil, "-"); actual != "" {
		t.Errorf("got %s; wanted empty string", actual)
	}
}
