package monitor

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

type postHandler struct{}

func (p *postHandler) Post(post *redditproto.Link) {}

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

func TestMonitors(t *testing.T) {
	mons, err := Monitors(
		&postHandler{},
		nil,
		&operator.MockOperator{},
		Forward,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(mons) != 0 {
		t.Errorf("got %d monitors; wanted 0", len(mons))
	}

	mons, err = Monitors(
		&mockInboxHandler{},
		[]string{"self"},
		&operator.MockOperator{},
		Backward,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(mons) != 4 {
		t.Fatalf("got %d monitors; wanted 1", len(mons))
	}

	mons, err = Monitors(
		&mockPostHandler{},
		[]string{"self"},
		&operator.MockOperator{},
		Backward,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(mons) != 1 {
		t.Fatalf("got %d monitors; wanted 1", len(mons))
	}
}

func TestBaseFromPath(t *testing.T) {
	mon, err := baseFromPath(
		&operator.MockOperator{
			ScrapeLinksReturn: []*redditproto.Link{
				&redditproto.Link{
					Name: stringPointer("name"),
				},
			},
		},
		"/r/self",
		nil,
		nil,
		nil,
		Forward,
	)
	if err != nil {
		t.Fatal(err)
	}

	b := mon.(*base)
	if b.dir != Forward {
		t.Errorf("got %d; wanted %d (Forward)", b.dir, Forward)
	}
	if b.handlePost != nil {
		t.Errorf("wanted post handler unset")
	}
	if b.path != "/r/self" {
		t.Errorf("got %s; wanted /r/self", b.path)
	}
}

func TestMerge(t *testing.T) {
	things := merge(
		[]*redditproto.Link{
			&redditproto.Link{
				CreatedUtc: float64Pointer(1),
				Name:       stringPointer("one"),
			},
			&redditproto.Link{
				CreatedUtc: float64Pointer(2),
				Name:       stringPointer("two"),
			},
		},
		[]*redditproto.Comment{
			&redditproto.Comment{
				CreatedUtc: float64Pointer(0),
				Name:       stringPointer("zero"),
			},
			&redditproto.Comment{
				CreatedUtc: float64Pointer(3),
				Name:       stringPointer("three"),
			},
		},
		[]*redditproto.Message{
			&redditproto.Message{
				CreatedUtc: float64Pointer(4),
				Name:       stringPointer("four"),
			},
			&redditproto.Message{
				CreatedUtc: float64Pointer(5),
				Name:       stringPointer("five"),
			},
		},
		Forward,
	)

	if len(things) != 6 {
		t.Fatalf("got %d things; wanted 6", len(things))
	}

	if things[0].GetName() != "zero" {
		t.Errorf("got %s; wanted zero", things[0].GetName())
	}

	if things[1].GetName() != "one" {
		t.Errorf("got %s; wanted one", things[1].GetName())
	}

	if things[2].GetName() != "two" {
		t.Errorf("got %s; wanted two", things[2].GetName())
	}

	if things[3].GetName() != "three" {
		t.Errorf("got %s; wanted three", things[3].GetName())
	}

	if things[4].GetName() != "four" {
		t.Errorf("got %s; wanted four", things[4].GetName())
	}

	if things[5].GetName() != "five" {
		t.Errorf("got %s; wanted five", things[5].GetName())
	}

	things = merge(
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
		Backward,
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
		&operator.MockOperator{
			IsThereThingErr: fmt.Errorf("an error"),
		},
	)
	if err == nil {
		t.Errorf("wanted error propagated from operator error")
	}

	broken, err = b.fixTip(
		&operator.MockOperator{
			IsThereThingReturn: false,
		},
	)
	if !broken {
		t.Errorf("got false; wanted true")
	}
}

func TestUpdateTip(t *testing.T) {
	b := &base{
		dir: Forward,
	}
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

	if b.tip[0] != "one" {
		t.Errorf("got %s; wanted one", b.tip[0])
	}

	if b.tip[1] != "two" {
		t.Errorf("got %s; wanted two", b.tip[1])
	}

	err = b.updateTip(
		nil,
		nil,
		nil,
		&operator.MockOperator{
			IsThereThingErr: fmt.Errorf("an error"),
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
	err := b.healthCheck(&operator.MockOperator{})
	if err != nil {
		t.Fatal(err)
	}
	if b.blanks != 1 {
		t.Errorf("got %d blanks; wanted 1", b.blanks)
	}

	b.blanks = b.blankThreshold
	err = b.healthCheck(
		&operator.MockOperator{
			IsThereThingReturn: true,
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
		&operator.MockOperator{
			IsThereThingErr: fmt.Errorf("an error"),
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
