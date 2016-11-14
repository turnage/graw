package lurker

import (
	"fmt"
	"strings"

	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/data"
	"github.com/turnage/graw/internal/reap"
)

// deletedAuthor is the author field of deleted posts on Reddit.
const deletedAuthor = "[deleted]"

// DoesNotExistErr indicates a value did not exist at an endpoint.
var DoesNotExistErr = fmt.Errorf("did not find expected values at endpoint")

// Lurker provides a high level interface for information fetching api calls to
// Reddit.
type Lurker interface {
	// Listing returns a harvest from a listing endpoint at Reddit.
	Listing(path, after string) (reap.Harvest, error)
	// Thread returns a Reddit post with a full parsed comment tree. The
	// permalink can be used as the path.
	Thread(path string) (*data.Post, error)
	// Exists returns whether a thing with the given name exists on Reddit
	// and is not deleted. A name is a type code (t#_) and an id, e.g.
	// "t1_fjsj3jf".
	Exists(name string) (bool, error)
}

type lurker struct {
	r reap.Reaper
}

func New(r reap.Reaper) Lurker {
	return &lurker{r: r}
}

func (l *lurker) Listing(path, after string) (reap.Harvest, error) {
	return l.r.Reap(
		path, api.WithDefaults(
			map[string]string{
				"limit":  "100",
				"before": after,
			},
		),
	)
}

func (l *lurker) Thread(path string) (*data.Post, error) {
	harvest, err := l.r.Reap(path+".json", api.WithDefaults(nil))
	if err != nil {
		return nil, err
	}

	if len(harvest.Posts) != 1 {
		return nil, DoesNotExistErr
	}

	return harvest.Posts[0], nil
}

func (l *lurker) Exists(name string) (bool, error) {
	path := "/api/info.json"

	// api/info doesn't provide message types; these need to be fetched from
	// a different url.
	if strings.HasPrefix(name, "t4_") {
		id := strings.TrimPrefix(name, "t4_")
		path = fmt.Sprintf("/message/messages/%s", id)
	}

	h, err := l.r.Reap(
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
