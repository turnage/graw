package reddit

// deletedAuthor is the author field of deleted posts on Reddit.
const deletedAuthor = "[deleted]"

// Scanner provides a high level interface for information fetching api calls to
// Reddit.
type Scanner interface {
	// Listing returns a harvest from a listing endpoint at Reddit.
	Listing(path, after string) (Harvest, error)
}

type scanner struct {
	r reaper
}

func newScanner(r reaper) Scanner {
	return &scanner{r: r}
}

func (s *scanner) Listing(path, after string) (Harvest, error) {
	return s.r.reap(
		path, map[string]string{
			"raw_json": "1",
			"limit":    "100",
			"before":   after,
		},
	)
}
