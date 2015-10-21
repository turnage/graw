package scanner

import (
	"fmt"
	"testing"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

func TestScanReportsScrapeErrors(t *testing.T) {
	expectedErr := fmt.Errorf("an error")
	sc := &listingScanner{
		tip:  []string{""},
		user: "user",
		op: &operator.MockOperator{
			UserContentErr: expectedErr,
		},
	}
	if _, _, err := sc.Scan(); err == nil {
		t.Errorf("wanted error for user content request failure")
	}

	sc = &listingScanner{
		tip:       []string{""},
		subreddit: "self",
		op: &operator.MockOperator{
			PostsErr: expectedErr,
		},
	}
	if _, _, err := sc.Scan(); err == nil {
		t.Errorf("wanted error for posts request failure")
	}
}

func TestScanReturnsOnlyNewThings(t *testing.T) {
	thingName := "name"
	sc := NewPostScanner(
		"self",
		&operator.MockOperator{
			PostsReturn: []*redditproto.Link{
				&redditproto.Link{Name: &thingName},
			},
		},
	)

	// The first scan should set the tip and not return any content.
	links, _, err := sc.Scan()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if len(links) != 0 {
		t.Fatalf("got %d links; wanted 0", len(links))
	}

	// After setting the tip, listings should be returned.
	links, _, err = sc.Scan()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("got %d links; wanted 1", len(links))
	}

	if links[0].GetName() != thingName {
		t.Errorf("got %s; wanted %s", links[0].GetName(), thingName)
	}
}

func TestScanReportsFixTipErrors(t *testing.T) {
	expectedErr := fmt.Errorf("an error")
	sc := NewUserScanner(
		"user",
		&operator.MockOperator{
			IsThereThingErr: expectedErr,
		},
	)
	sc.blanks = sc.blankThreshold + 1
	if _, _, err := sc.Scan(); err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}
}

func TestScanIncreasesBlankThreshold(t *testing.T) {
	sc := NewUserScanner(
		"user",
		&operator.MockOperator{
			IsThereThingReturn: true,
		},
	)
	sc.blanks = sc.blankThreshold + 1
	_, comments, err := sc.Scan()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("got %v; wanted nothing", comments)
	}
	if sc.blankThreshold <= defaultBlankThreshold {
		t.Errorf(
			"got %d; wanted higher than %d",
			sc.blankThreshold,
			defaultBlankThreshold)
	}
}

func TestFetchTip(t *testing.T) {
	sc := &listingScanner{
		subreddit: "self",
		tip:       []string{""},
	}

	sc.op = &operator.MockOperator{
		PostsErr: fmt.Errorf("an error"),
	}
	if _, _, err := sc.fetchTip(); err == nil {
		t.Errorf("wanted error for request failure")
	}

	sc.tip = make([]string, maxTipSize)
	for i := 0; i < maxTipSize; i++ {
		sc.tip = append(sc.tip, "id")
	}
	linkName := "anything"
	sc.op = &operator.MockOperator{
		PostsErr: nil,
		PostsReturn: []*redditproto.Link{
			&redditproto.Link{Name: &linkName},
		},
	}

	links, _, err := sc.fetchTip()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if sc.tip[len(sc.tip)-1] != linkName {
		t.Errorf(
			"got tip %s; wanted %s",
			sc.tip[len(sc.tip)-1],
			linkName)
	}

	if len(links) != 1 {
		t.Fatalf("got %d links; expected 1", len(links))
	}

	if links[0].GetName() != linkName {
		t.Errorf(
			"got thread name %s; wanted %s",
			links[0].GetName(),
			linkName)
	}

	sc.tip = []string{""}
	links, _, err = sc.fetchTip()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if links != nil {
		t.Errorf("got %v; wanted no links for adjustment round", links)
	}
}

func TestFixTip(t *testing.T) {
	sc := &listingScanner{
		tip: []string{"1", "2", "3"},
	}

	sc.op = &operator.MockOperator{
		IsThereThingErr: fmt.Errorf("an error"),
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

	sc.op = &operator.MockOperator{
		IsThereThingReturn: true,
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
	sc := &listingScanner{
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
