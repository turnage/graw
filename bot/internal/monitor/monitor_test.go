package monitor

import (
	"reflect"
	"strings"
	"testing"
)

func TestMonitorToggles(t *testing.T) {
	monitor := &Monitor{
		monitoredSubreddits: make(map[string]bool),
		monitoredThreads:    make(map[string]bool),
	}

	monitor.MonitorSubreddits("awww", "self")
	if !strings.Contains(
		monitor.subredditQuery,
		"aww",
	) || !strings.Contains(
		monitor.subredditQuery,
		"self",
	) {
		t.Errorf(
			"got %s; wanted awww+self (any order of)",
			monitor.subredditQuery)
	}
	monitor.UnmonitorSubreddits("awww")
	if monitor.subredditQuery != "self" {
		t.Errorf("got %s; wanted self", monitor.subredditQuery)
	}

	monitor.MonitorThreads("harry", "potter")
	if !strings.Contains(
		monitor.threadQuery,
		"harry",
	) || !strings.Contains(
		monitor.threadQuery,
		"potter",
	) {
		t.Errorf(
			"got %s; wanted harry,potter (any order of)",
			monitor.threadQuery)
	}
	monitor.UnmonitorThreads("potter")
	if monitor.threadQuery != "harry" {
		t.Errorf("got %s; wanted harry", monitor.threadQuery)
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
