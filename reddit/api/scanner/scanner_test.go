package scanner

import (
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/pretty"

	"github.com/turnage/graw/reddit"
	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/reap"
)

func TestListing(t *testing.T) {
	h := reap.Harvest{
		Comments: []*reddit.Comment{
			&reddit.Comment{
				Body: "text",
			},
		},
	}
	s := New(api.ReaperWhich(h, nil))

	actual, err := s.Listing("/messages", "")
	if err != nil {
		t.Errorf("error lurking listing: %v", err)
	}

	if diff := pretty.Compare(&h, &actual); diff != "" {
		t.Errorf("harvest unexpected; diff: %s", diff)
	}
}

func TestExists(t *testing.T) {
	empty := reap.Harvest{}
	post := reap.Harvest{Posts: []*reddit.Post{&reddit.Post{}}}
	comment := reap.Harvest{Comments: []*reddit.Comment{&reddit.Comment{}}}
	message := reap.Harvest{Messages: []*reddit.Message{&reddit.Message{}}}
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
		r := api.ReaperWhich(test.h, test.err)
		s := New(r)
		exists, err := s.Exists(test.input)
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

		if r.Path != test.path {
			t.Errorf("got path %s; wanted %s", r.Path, test.path)
		}
	}
}
