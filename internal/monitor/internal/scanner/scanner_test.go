package scanner

import (
	"fmt"
	"testing"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

func TestScanReportsScrapeErrors(t *testing.T) {
	sc := &Scanner{
		tip: []string{""},
	}

	expectedErr := fmt.Errorf("an error")
	sc.op = &operator.MockOperator{
		ScrapeErr: expectedErr,
	}
	if _, err := sc.Scan(); err == nil {
		t.Errorf("wanted error for request failure")
	}
}

func TestScanReturnsOnlyNewThings(t *testing.T) {
	thingName := "name"
	sc := New(
		"",
		&operator.MockOperator{
			ScrapeReturn: []operator.Thing{
				&redditproto.Link{Name: &thingName},
			},
		},
	)

	// The first scan should set the tip and not return any Things.
	things, err := sc.Scan()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if len(things) != 0 {
		t.Fatalf("got %d things; wanted 0", len(things))
	}

	// After setting the tip, listings should be returned.
	things, err = sc.Scan()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if len(things) != 1 {
		t.Fatalf("got %d things; wanted 1", len(things))
	}

	if things[0].GetName() != thingName {
		t.Errorf("got %s; wanted %s", things[0].GetName(), thingName)
	}
}

func TestScanReportsFixTipErrors(t *testing.T) {
	expectedErr := fmt.Errorf("an error")
	sc := New(
		"",
		&operator.MockOperator{
			ThreadsErr: expectedErr,
		},
	)
	sc.blanks = sc.blankThreshold + 1
	if _, err := sc.Scan(); err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}
}

func TestScanIncreasesBlankThreshold(t *testing.T) {
	sc := New(
		"",
		&operator.MockOperator{
			ThreadsReturn: []*redditproto.Link{
				&redditproto.Link{},
			},
		},
	)
	sc.blanks = sc.blankThreshold + 1
	things, err := sc.Scan()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if len(things) != 0 {
		t.Errorf("got %v; wanted nothing", things)
	}
	if sc.blankThreshold <= defaultBlankThreshold {
		t.Errorf(
			"got %d; wanted higher than %d",
			sc.blankThreshold,
			defaultBlankThreshold)
	}
}

func TestFetchTip(t *testing.T) {
	sc := &Scanner{
		tip: []string{""},
	}

	sc.op = &operator.MockOperator{
		ScrapeErr: fmt.Errorf("an error"),
	}
	if _, err := sc.fetchTip(); err == nil {
		t.Errorf("wanted error for request failure")
	}

	sc.tip = make([]string, maxTipSize)
	for i := 0; i < maxTipSize; i++ {
		sc.tip = append(sc.tip, "id")
	}
	thingName := "anything"
	sc.op = &operator.MockOperator{
		ScrapeErr: nil,
		ScrapeReturn: []operator.Thing{
			&redditproto.Link{Name: &thingName},
		},
	}

	things, err := sc.fetchTip()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if sc.tip[len(sc.tip)-1] != thingName {
		t.Errorf(
			"got tip %s; wanted %s",
			sc.tip[len(sc.tip)-1],
			thingName)
	}

	if len(things) != 1 {
		t.Fatalf("got %d things; expected 1", len(things))
	}

	if things[0].GetName() != thingName {
		t.Errorf(
			"got thread name %s; wanted %s",
			things[0].GetName(),
			thingName)
	}

	sc.tip = []string{""}
	things, err = sc.fetchTip()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if things != nil {
		t.Errorf("got %v; wanted no things for adjustment round", things)
	}
}

func TestFixTip(t *testing.T) {
	sc := &Scanner{
		tip: []string{"1", "2", "3"},
	}

	sc.op = &operator.MockOperator{
		ThreadsErr: fmt.Errorf("an error"),
	}
	if _, err := sc.fixTip(); err == nil {
		t.Errorf("wanted error for request failure")
	}

	sc.op = &operator.MockOperator{}
	shaved, err := sc.fixTip()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if !shaved {
		t.Errorf("wanted indication that the tip has been shaved")
	}

	if sc.tip[len(sc.tip)-1] != "2" {
		t.Errorf(
			"got %s; wanted tip shaved to 2",
			sc.tip[len(sc.tip)-1])
	}

	nameTwo := "2"
	sc.op = &operator.MockOperator{
		ThreadsReturn: []*redditproto.Link{
			&redditproto.Link{Name: &nameTwo},
		},
	}
	shaved, err = sc.fixTip()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if shaved {
		t.Errorf("tip was shaved when it should not have been")
	}
}

func TestShaveTip(t *testing.T) {
	sc := &Scanner{
		tip: []string{"1", "2"},
	}

	sc.shaveTip()
	if sc.tip[len(sc.tip)-1] != "1" {
		t.Errorf(
			"got %s; wanted tip shaved to 1",
			sc.tip[len(sc.tip)-1])
	}

	sc.shaveTip()
	if len(sc.tip) != 1 {
		t.Errorf("tip is %d long; wanted 1 blank tip", len(sc.tip))
	}

	if sc.tip[0] != "" {
		t.Errorf("got %s; wanted empty string", sc.tip[0])
	}
}
