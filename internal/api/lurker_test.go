package api

import (
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/pretty"

	"github.com/turnage/graw/internal/data"
	"github.com/turnage/graw/internal/reap"
)

func TestListing(t *testing.T) {
	h := reap.Harvest{
		Comments: []*data.Comment{
			&data.Comment{
				Body: "text",
			},
		},
	}
	l := NewLurker(reaperWhich(h, nil))

	actual, err := l.Listing("/messages", "")
	if err != nil {
		t.Errorf("error lurking listing: %v", err)
	}

	if diff := pretty.Compare(&h, &actual); diff != "" {
		t.Errorf("harvest unexpected; diff: %s", diff)
	}
}

func TestThread(t *testing.T) {
	h := reap.Harvest{
		Posts: []*data.Post{
			&data.Post{
				SelfText: "hello",
			},
		},
	}
	l := NewLurker(reaperWhich(h, nil))

	post, err := l.Thread("")
	if err != nil {
		t.Errorf("error pulling thread: %v", err)
	}

	if diff := pretty.Compare(post, h.Posts[0]); diff != "" {
		t.Errorf("post incorrect; diff: %s", diff)
	}
}

func TestThreadReturnsEmpty(t *testing.T) {
	l := NewLurker(reaperWhich(reap.Harvest{}, nil))
	_, err := l.Thread("")
	if err != DoesNotExistErr {
		t.Errorf("err unexpected; wanted DoesNotExistErr; got %v", err)
	}
}

func TestExists(t *testing.T) {
	empty := reap.Harvest{}
	post := reap.Harvest{Posts: []*data.Post{&data.Post{}}}
	comment := reap.Harvest{Comments: []*data.Comment{&data.Comment{}}}
	message := reap.Harvest{Messages: []*data.Message{&data.Message{}}}
	fail := fmt.Errorf("a failure")

	for _, test := range []struct {
		input string
		h     reap.Harvest
		err   error

		exists bool
		path   string
	}{
		{"t1_ffjdkdf", empty, nil, false, "/api/info.json"},
		{"t4_ffjdkdf", empty, nil, false, "/message/messages/ffjdkdf"},
		{"t2_fffjsdj", post, nil, true, "/api/info.json"},
		{"t4_fffjsdj", message, nil, true, "/message/messages/fffjsdj"},
		{"t1_fffjsdj", comment, nil, true, "/api/info.json"},
		{"t1_fffjsdj", comment, fail, false, "/api/info.json"},
	} {
		r := reaperWhich(test.h, test.err)
		l := NewLurker(r)
		exists, err := l.Exists(test.input)
		if err != test.err {
			t.Errorf("got err %v; wanted %v", err, test.err)
		}

		if exists != test.exists {
			t.Errorf(
				"got existence result %v; wanted %v",
				exists,
				test.exists,
			)
		}

		if r.path != test.path {
			t.Errorf("got path %s; wanted %s", r.path, test.path)
		}
	}
}
