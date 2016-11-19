package reddit

import (
	"fmt"
)

// PostDoesNotExistErr indicates a post does not exist.
var PostDoesNotExistErr = fmt.Errorf("The requested post does not exist.")

// Lurker provides a high level interface for fetching information from Reddit.
type Lurker interface {
	// Thread returns a Reddit post with a full parsed comment tree. The
	// permalink can be used as the path.
	Thread(permalink string) (*Post, error)
}

type lurker struct {
	r reaper
}

func newLurker(r reaper) Lurker {
	return &lurker{r: r}
}

func (s *lurker) Thread(permalink string) (*Post, error) {
	harvest, err := s.r.reap(permalink+".json", withDefaultAPIArgs(nil))
	if err != nil {
		return nil, err
	}

	if len(harvest.Posts) != 1 {
		return nil, PostDoesNotExistErr
	}

	return harvest.Posts[0], nil
}
