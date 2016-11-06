package rsort

import (
	"testing"

	"github.com/turnage/graw/internal/data"
	"github.com/turnage/graw/internal/reap"
)

func TestSort(t *testing.T) {
	names := New().Sort(
		reap.Harvest{
			Posts: []*data.Post{
				&data.Post{CreatedUTC: 1, Name: "1"},
				&data.Post{CreatedUTC: 7, Name: "7"},
				&data.Post{CreatedUTC: 2, Name: "2"},
			},
			Comments: []*data.Comment{
				&data.Comment{CreatedUTC: 5, Name: "5"},
				&data.Comment{CreatedUTC: 0, Name: "0"},
			},
			Messages: []*data.Message{
				&data.Message{CreatedUTC: 6, Name: "6"},
			},
		},
	)

	if len(names) != 6 {
		t.Errorf("unexpected length; got %d; wanted %d", len(names), 6)
	}

	// Younger elements (those with later/higher creation times) should come
	// first.
	for i, name := range []string{"7", "6", "5", "2", "1", "0"} {
		if names[i] != name {
			t.Errorf("%d wrong; got %s vs %v", i, names[i], name)
		}
	}
}
