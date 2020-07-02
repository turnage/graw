package reddit

import (
	"fmt"
	"strings"
)

// Lurker defines browsing behavior.
type Lurker interface {
	// Thread returns a Reddit post with a fully parsed comment tree.
	Thread(permalink string) (*Post, error)
	// MoreChildren returns a complete comment tree for a Reddit post
	MoreChildren(more *More, link string) ([]*Comment, error)
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

func (s *lurker) MoreChildren(more *More, link string) ([]*Comment, error) {
	reaperParams := map[string]string{
		"api_type": "json",
		"link_id":  link,
		"children": strings.Join(more.Children, ","),
	}

	harvest, err := s.r.reap("/api/morechildren", reaperParams)
	if err != nil {
		return nil, err
	}

	comments := []*Comment{}

	for _, m := range harvest.Mores {
		cs, err := s.MoreChildren(m, link)
		if err != nil {
			return nil, err
		}
		comments = append(comments, cs...)
	}

	fmt.Println("Comments:", len(comments), len(harvest.Comments))

	return append(comments, harvest.Comments...), nil

}
