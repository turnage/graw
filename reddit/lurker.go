package reddit

// Lurker defines browsing behavior.
type Lurker interface {
	// Thread returns a Reddit post with a fully parsed comment tree.
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
