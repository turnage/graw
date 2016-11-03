package api

import (
	"fmt"

	"github.com/turnage/graw/internal/data"
	"github.com/turnage/graw/internal/reap"
)

var (
	EmptyHarvestErr = fmt.Errorf("did not find expected values at endpoint")
)

type Lurker interface {
	Listing(subreddit, after string) (reap.Harvest, error)
	Thread(path string) (*data.Post, error)
}

type lurker struct {
	r reap.Reaper
}

func NewLurker(r reap.Reaper) Lurker {
	return &lurker{r: r}
}

func (l *lurker) Listing(subreddit, after string) (reap.Harvest, error) {
	return l.r.Reap(
		"/r/"+subreddit,
		withDefaults(map[string]string{"limit": "100"}),
	)
}

func (l *lurker) Thread(path string) (*data.Post, error) {
	harvest, err := l.r.Reap(path+".json", withDefaults(nil))
	if err != nil {
		return nil, err
	}

	if len(harvest.Posts) != 1 {
		return nil, EmptyHarvestErr
	}

	return harvest.Posts[0], nil
}
