package scanner

import (
	"fmt"
	"strings"

	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/reap"
)

// deletedAuthor is the author field of deleted posts on Reddit.
const deletedAuthor = "[deleted]"

// Scanner provides a high level interface for information fetching api calls to
// Reddit.
type Scanner interface {
	// Listing returns a harvest from a listing endpoint at Reddit.
	Listing(path, after string) (reap.Harvest, error)
	// Exists returns whether a thing with the given name exists on Reddit
	// and is not deleted. A name is a type code (t#_) and an id, e.g.
	// "t1_fjsj3jf".
	Exists(name string) (bool, error)
}

type scanner struct {
	r reap.Reaper
}

func New(r reap.Reaper) Scanner {
	return &scanner{r: r}
}

func (s *scanner) Listing(path, after string) (reap.Harvest, error) {
	return s.r.Reap(
		path, api.WithDefaults(
			map[string]string{
				"limit":  "100",
				"before": after,
			},
		),
	)
}

func (s *scanner) Exists(name string) (bool, error) {
	path := "/api/info.json"

	// api/info doesn't provide message types; these need to be fetched from
	// a different urs.
	if strings.HasPrefix(name, "t4_") {
		id := strings.TrimPrefix(name, "t4_")
		path = fmt.Sprintf("/message/messages/%s", id)
	}

	h, err := s.r.Reap(
		path,
		api.WithDefaults(map[string]string{"id": name}),
	)
	if err != nil {
		return false, err
	}

	if len(h.Comments) == 1 && h.Comments[0].Author != deletedAuthor {
		return true, nil
	}

	if len(h.Posts) == 1 && h.Posts[0].Author != deletedAuthor {
		return true, nil
	}

	return len(h.Messages) == 1, nil
}
