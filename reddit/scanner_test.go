package reddit

import (
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestListing(t *testing.T) {
	h := Harvest{
		Comments: []*Comment{
			&Comment{
				Body: "text",
			},
		},
	}
	s := newScanner(reaperWhich(h, nil))

	actual, err := s.Listing("/messages", "")
	if err != nil {
		t.Errorf("error lurking listing: %v", err)
	}

	if diff := pretty.Compare(&h, &actual); diff != "" {
		t.Errorf("harvest unexpected; diff: %s", diff)
	}
}

func TestExists(t *testing.T) {
	empty := Harvest{}
	post := Harvest{Posts: []*Post{&Post{}}}
	comment := Harvest{Comments: []*Comment{&Comment{}}}
	message := Harvest{Messages: []*Message{&Message{}}}
	fail := fmt.Errorf("a failure")

	for _, test := range []struct {
		input string
		h     Harvest
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
		s := newScanner(r)
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

		if r.path != test.path {
			t.Errorf("got path %s; wanted %s", r.path, test.path)
		}
	}
}
