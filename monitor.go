package graw

import (
	"fmt"
	"strings"
	"time"

	"github.com/turnage/redditproto"
)

// subredditMonitor monitors subreddits for new posts, and feeds the posts it
// finds through the posts channel.
type subredditMonitor struct {
	// posts is the channel new posts are fed through.
	Posts chan *redditproto.Link
	// errors is the channel errors are fed through before the
	// subredditMonitor stops, so its controller knows why it failed.
	Errors chan error
	// subreddits is the list of subreddits the monitor monitors.
	Subreddits []string
	// kill is a channel the subredditMonitor's controller can use to kill
	// it.
	Kill chan bool
	// RefreshRate is the amount of times per minute the monitor will check
	// for new posts.
	RefreshRate int

	// last is the fullname of the freshest post at the last check
	last string
	// lastURL is the url of the freshest post at the last check
	lastURL string
}

// Run continuously polls monitored subreddits for new posts.
func (s *subredditMonitor) Run(cli client) {
	_, err := s.tip(cli, 1)
	if err != nil {
		s.Errors <- err
		return
	}

	for true {
		select {
		case <-time.After(time.Minute / time.Duration(s.RefreshRate)):
			posts, err := s.tip(cli, 100)
			if err != nil {
				s.Errors <- err
				return
			}
			fmt.Printf("Found %d posts since %s.\n", len(posts), s.lastURL)
			for _, post := range posts {
				s.Posts <- post
			}
		case <-s.Kill:
			return
		}
	}
}

// tip returns the posts made since the last check, from the previous tip up to
// lim.
func (s *subredditMonitor) tip(cli client, lim int) ([]*redditproto.Link, error) {
	posts, err := scrape(
		cli,
		strings.Join(s.Subreddits, "+"),
		"new",
		"",
		s.last,
		lim)
	if err != nil {
		return nil, err
	}

	if len(posts) > 0 {
		s.last = posts[len(posts)-1].GetName()
		s.lastURL = posts[len(posts)-1].GetPermalink()
	}

	return posts, nil
}
