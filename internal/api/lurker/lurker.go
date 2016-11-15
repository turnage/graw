package lurker

import (
	"fmt"

	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/data"
	"github.com/turnage/graw/internal/reap"
)

// DoesNotExistErr indicates a value did not exist at an endpoint.
var DoesNotExistErr = fmt.Errorf("did not find expected values at endpoint")

// Lurker provides a high level interface for fetching information from Reddit.
type Lurker interface {
	// Thread returns a Reddit post with a full parsed comment tree. The
	// permalink can be used as the path.
	Thread(permalink string) (*data.Post, error)
}

type lurker struct {
	r reap.Reaper
}

func New(r reap.Reaper) Lurker {
	return &lurker{r: r}
}

func (s *lurker) Thread(permalink string) (*data.Post, error) {
	harvest, err := s.r.Reap(permalink+".json", api.WithDefaults(nil))
	if err != nil {
		return nil, err
	}

	if len(harvest.Posts) != 1 {
		return nil, DoesNotExistErr
	}

	return harvest.Posts[0], nil
}
