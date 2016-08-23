package monitor

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/turnage/redditproto"
)

type handler struct {
	calls        int
	postCalls    int
	commentCalls int
	messageCalls int
}

func (h *handler) post(post *redditproto.Link) {
	h.postCalls++
	h.calls++
}

func (h *handler) comment(comment *redditproto.Comment) {
	h.commentCalls++
	h.calls++
}

func (h *handler) message(message *redditproto.Message) {
	h.messageCalls++
	h.calls++
}

func float64Pointer(val float64) *float64 {
	return &val
}

func stringPointer(val string) *string {
	return &val
}

func MockScraper(
	links []*redditproto.Link,
	comments []*redditproto.Comment,
	messages []*redditproto.Message,
	err error,
) Scraper {
	return func(id, tip string, limit int) (
		[]*redditproto.Link,
		[]*redditproto.Comment,
		[]*redditproto.Message,
		error,
	) {
		return links, comments, messages, err
	}
}

func TestBaseFromPath(t *testing.T) {
	han := &handler{}
	mon, err := baseFromPath(
		MockScraper(
			[]*redditproto.Link{
				&redditproto.Link{
					Name: stringPointer("name"),
				},
			},
			nil, nil, nil,
		),
		"/r/self",
		nil,
		han.comment,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	b := mon.(*base)
	if b.handleComment == nil {
		t.Errorf("wanted comment handler set")
	}
	if b.path != "/r/self" {
		t.Errorf("got %s; wanted /r/self", b.path)
	}
	mon, err = baseFromPath(
		MockScraper(
			[]*redditproto.Link{
				&redditproto.Link{
					Name: stringPointer("name"),
				},
			}, nil, nil, nil,
		),
		"/r/self",
		nil,
		nil,
		nil,
	)
	if err == nil {
		t.Errorf("wanted error if no handlers are provided")
	}
}

func TestMerge(t *testing.T) {
	things := merge(
		[]*redditproto.Link{
			&redditproto.Link{
				CreatedUtc: float64Pointer(2),
				Name:       stringPointer("two"),
			},
			&redditproto.Link{
				CreatedUtc: float64Pointer(1),
				Name:       stringPointer("one"),
			},
		},
		[]*redditproto.Comment{
			&redditproto.Comment{
				CreatedUtc: float64Pointer(3),
				Name:       stringPointer("three"),
			},
			&redditproto.Comment{
				CreatedUtc: float64Pointer(0),
				Name:       stringPointer("zero"),
			},
		},
		[]*redditproto.Message{
			&redditproto.Message{
				CreatedUtc: float64Pointer(5),
				Name:       stringPointer("five"),
			},
			&redditproto.Message{
				CreatedUtc: float64Pointer(4),
				Name:       stringPointer("four"),
			},
		},
	)

	if len(things) != 6 {
		t.Fatalf("got %d things; wanted 6", len(things))
	}

	if things[0].GetName() != "five" {
		t.Errorf("got %s; wanted five", things[0].GetName())
	}

	if things[1].GetName() != "four" {
		t.Errorf("got %s; wanted four", things[1].GetName())
	}

	if things[2].GetName() != "three" {
		t.Errorf("got %s; wanted three", things[2].GetName())
	}

	if things[3].GetName() != "two" {
		t.Errorf("got %s; wanted two", things[3].GetName())
	}

	if things[4].GetName() != "one" {
		t.Errorf("got %s; wanted one", things[4].GetName())
	}

	if things[5].GetName() != "zero" {
		t.Errorf("got %s; wanted zero", things[5].GetName())
	}
}

func TestShaveTip(t *testing.T) {
	b := &base{
		tip: []string{"1", "2"},
	}
	b.shaveTip()
	if !reflect.DeepEqual(b.tip, []string{"2"}) {
		t.Errorf("got %v\n; wanted %v", b.tip, []string{"2"})
	}

	b = &base{
		tip: nil,
	}
	b.shaveTip()
	if !reflect.DeepEqual(b.tip, []string{""}) {
		t.Errorf("got %v\n; wanted %v", b.tip, []string{""})
	}
}

func TestFixTip(t *testing.T) {
	b := &base{
		tip: []string{"1"},
	}

	broken, err := b.fixTip(
		func(id string) (bool, error) {
			return false, fmt.Errorf("an error")
		},
	)
	if err == nil {
		t.Errorf("wanted error propagated from operator error")
	}

	broken, err = b.fixTip(
		func(id string) (bool, error) {
			return false, nil
		},
	)
	if !broken {
		t.Errorf("got false; wanted true")
	}
}

func TestUpdateTip(t *testing.T) {
	b := &base{}
	for i := 0; i < maxTipSize; i++ {
		b.tip = append(b.tip, strconv.Itoa(i))
	}

	err := b.updateTip(
		[]*redditproto.Link{
			&redditproto.Link{
				CreatedUtc: float64Pointer(2),
				Name:       stringPointer("two"),
			},
			&redditproto.Link{
				CreatedUtc: float64Pointer(1),
				Name:       stringPointer("one"),
			},
		},
		nil,
		nil,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(b.tip) != maxTipSize {
		t.Fatalf("got size %d tip log; wanted %d", len(b.tip), maxTipSize)
	}

	if b.tip[0] != "two" {
		t.Errorf("got %s; wanted two", b.tip[0])
	}

	if b.tip[1] != "one" {
		t.Errorf("got %s; wanted one", b.tip[1])
	}

	err = b.updateTip(
		nil,
		nil,
		nil,
		func(id string) (bool, error) {
			return false, fmt.Errorf("an error")
		},
	)
	if err == nil {
		t.Errorf("wanted error propagated from healthCheck")
	}
}

func TestHealthCheck(t *testing.T) {
	b := &base{
		blankThreshold: blankThreshold,
		tip:            []string{""},
	}
	err := b.healthCheck(
		func(id string) (bool, error) {
			return false, nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if b.blanks != 1 {
		t.Errorf("got %d blanks; wanted 1", b.blanks)
	}

	b.blanks = b.blankThreshold
	err = b.healthCheck(
		func(id string) (bool, error) {
			return true, nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if b.blanks != 0 {
		t.Errorf("got %d blanks; wanted 0", b.blanks)
	}
	if b.blankThreshold <= blankThreshold {
		t.Errorf("got %d; wanted > %d", b.blankThreshold, blankThreshold)
	}

	b.blanks = b.blankThreshold
	err = b.healthCheck(
		func(id string) (bool, error) {
			return false, fmt.Errorf("an error")
		},
	)
	if err == nil {
		t.Fatalf("wanted error propagated from operator")
	}
}

func TestDispatch(t *testing.T) {
	hand := &handler{}
	b := &base{
		handlePost:    hand.post,
		handleComment: hand.comment,
		handleMessage: hand.message,
	}
	b.dispatch(
		[]*redditproto.Link{
			&redditproto.Link{},
		},
		[]*redditproto.Comment{
			&redditproto.Comment{},
		},
		[]*redditproto.Message{
			&redditproto.Message{},
		},
	)

	for i := 0; i < 100 && hand.calls < 3; i++ {
		time.Sleep(10 * time.Millisecond)
	}

	if hand.postCalls != 1 {
		t.Errorf("got %d post calls; wanted 1", hand.postCalls)
	}

	if hand.commentCalls != 1 {
		t.Errorf("got %d comment calls; wanted 1", hand.commentCalls)
	}

	if hand.messageCalls != 1 {
		t.Errorf("got %d message calls; wanted 1", hand.messageCalls)
	}
}
