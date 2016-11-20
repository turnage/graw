package reddit

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
	harvest, err := s.r.reap(
		permalink+".json",
		map[string]string{"raw_json": "1"},
	)
	if err != nil {
		return nil, err
	}

	if len(harvest.Posts) != 1 {
		return nil, ThreadDoesNotExistErr
	}

	return harvest.Posts[0], nil
}
