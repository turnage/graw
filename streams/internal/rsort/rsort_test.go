package rsort

import (
	"testing"

	"github.com/turnage/graw/reddit"
)

func TestSort(t *testing.T) {
	names := New().Sort(
		reddit.Harvest{
			Posts: []*reddit.Post{
				{CreatedUTC: 1, Name: "1"},
				{CreatedUTC: 7, Name: "7"},
				{CreatedUTC: 2, Name: "2"},
			},
			Comments: []*reddit.Comment{
				{CreatedUTC: 5, Name: "5"},
				{CreatedUTC: 0, Name: "0"},
			},
			Messages: []*reddit.Message{
				{CreatedUTC: 6, Name: "6"},
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
