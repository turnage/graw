package lurker

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"

	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/data"
	"github.com/turnage/graw/internal/reap"
)

func TestThread(t *testing.T) {
	h := reap.Harvest{
		Posts: []*data.Post{
			&data.Post{
				SelfText: "hello",
			},
		},
	}
	s := New(api.ReaperWhich(h, nil))

	post, err := s.Thread("")
	if err != nil {
		t.Errorf("error pulling thread: %v", err)
	}

	if diff := pretty.Compare(post, h.Posts[0]); diff != "" {
		t.Errorf("post incorrect; diff: %s", diff)
	}
}

func TestThreadReturnsEmpty(t *testing.T) {
	s := New(api.ReaperWhich(reap.Harvest{}, nil))
	_, err := s.Thread("")
	if err != DoesNotExistErr {
		t.Errorf("err unexpected; wanted DoesNotExistErr; got %v", err)
	}
}
