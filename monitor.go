package graw

import (
	"time"
	"strings"

	"github.com/turnage/redditproto"
)

// subredditMonitor monitors subreddits for new posts, and feeds the posts it
// finds through the posts channel.
type subredditMonitor struct {
	// cli is used for executing network requests to reddit.
	cli client
	// posts is the channel new posts are fed through.
	posts chan *redditproto.Link
	// errors is the channel errors are fed through before the
	// subredditMonitor stops, so its controller knows why it failed.
	errors chan error
	// subreddits is the list of subreddits the monitor monitors.
	subreddits []string
	// kill is a channel the subredditMonitor's controller can use to kill
	// it.
	kill chan bool
	// last is the fullname of the freshest post as the last check
	last string
}

// run continuously polls monitored subreddits for new posts.
func (s *subredditMonitor) run() {
	_, err := s.tip(1)
	if err != nil {
		s.errors <- err
		return
	}

	for true {
		select {
		case <- time.After(3 * time.Second):
			posts, err := s.tip(100)
			if err != nil {
				s.errors <- err
				return
			}
			for _, post := range posts {
				s.posts <- post
			}
		case <- s.kill:
			return
		}
	}
}

// tip returns the posts made since the last check, from the previous tip up to
// lim.
func (s *subredditMonitor) tip(lim int) ([]*redditproto.Link, error) {
	posts, err := scrape(
		s.cli,
		strings.Join(s.subreddits, "+"),
		"new",
		"",
		s.last,
		lim)
	if err != nil {
		return nil, err
	}

	if len(posts) > 0 {
		s.last = posts[len(posts) - 1].GetName()
	}

	return posts, nil
}
