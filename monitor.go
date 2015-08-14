package graw

import (
	"container/list"
	"strings"
	"time"

	"github.com/turnage/redditproto"
)

const (
	// The amount of fallback threads to remember in case a reference
	// thread is deleted while it is in use.
	fallbackCount = 10
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
	// RefreshRate is the amount of times per minute the monitor will check
	// for new posts.
	RefreshRate int

	// last holds fullnames of the freshest posts at the last check. These
	// are used to differentiate new from old posts.
	last *list.List
}

// Run continuously polls monitored subreddits for new posts.
func (s *subredditMonitor) Run(cli client) {
	s.last = list.New()
	s.last.PushFront("")
	skipUpdate := false

	_, err := s.tip(cli, fallbackCount)
	if err != nil {
		s.Errors <- err
		return
	}

	for true {
		time.Sleep(time.Minute / time.Duration(s.RefreshRate))
		if skipUpdate {
			skipUpdate = false
			continue
		}

		posts, err := s.tip(cli, 100)
		if err != nil {
			s.Errors <- err
			return
		}

		if len(posts) == 0 {
			skipUpdate = true
			valid, err := s.validTip(cli)
			if err != nil {
				s.Errors <- err
				return
			}
			if !valid {
				s.last.Remove(s.last.Front())
				if s.last.Len() < 1 {
					s.last.PushFront("")
				}
			}
		} else {
			for _, post := range posts {
				s.Posts <- post
			}
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
		s.last.Front().Value.(string),
		lim)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		s.last.PushFront(posts[len(posts)-i-1].GetName())
		if s.last.Len() > fallbackCount {
			toRemove := s.last.Len() - fallbackCount
			for i := 0; i < toRemove; i++ {
				s.last.Remove(s.last.Back())
			}
		}
	}

	return posts, nil
}

// validTip checks that the post the monitor is using as the "latest" (a
// reference point from which to choose new threads) is still valid. If it has
// been deleted, requesting "newer" posts than the deleted thread will not work,
// and the monitor will think there are no new posts.
func (s *subredditMonitor) validTip(cli client) (bool, error) {
	link, err := threads(cli, s.last.Front().Value.(string))
	if err != nil {
		return false, err
	}

	if len(link) == 1 && link[0].GetAuthor() != "[deleted]" {
		return true, nil
	}

	return false, nil
}
