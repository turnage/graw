package reddit

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestThread(t *testing.T) {
	h := Harvest{
		Posts: []*Post{
			&Post{
				SelfText: "hello",
			},
		},
	}
	s := newLurker(reaperWhich(h, nil))

	post, err := s.Thread("")
	if err != nil {
		t.Errorf("error pulling thread: %v", err)
	}

	if diff := pretty.Compare(post, h.Posts[0]); diff != "" {
		t.Errorf("post incorrect; diff: %s", diff)
	}
}

func TestThreadReturnsEmpty(t *testing.T) {
	s := newLurker(reaperWhich(Harvest{}, nil))
	_, err := s.Thread("")
	if err != PostDoesNotExistErr {
		t.Errorf("err unexpected; wanted DoesNotExistErr; got %v", err)
	}
}
