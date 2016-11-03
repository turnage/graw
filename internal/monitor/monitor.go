// Package monitor includes monitors for different parts of Reddit, such as a
// user inbox or a subreddit's post feed.
package monitor

import (
	"fmt"
	"sort"

	"github.com/turnage/graw/internal/data"
)

const (
	// The blank threshold is the amount of updates returning 0 new
	// elements in the monitored listing the monitor will tolerate before
	// suspecting the tip of the listing has been deleted.
	blankThreshold = 2
	// maxTipSize is the maximum size of the tip log (number of backup tips
	// + the current tip).
	maxTipSize = 10
)

// PostHandler handles new posts in watched subreddits.
type PostHandler interface {
	// Post handles a new post in a watched subreddit.
	Post(p *data.Post) error
}

// Monitor defines the controls for a Monitor.
type Monitor interface {
	// Update will check for new events, and send them to the Monitor's
	// handlers.
	Update() error
}

// monitor describes the core of any monitor; most monitors will use its fields and
// methods.
type monitor struct {
	// blanks is the number of. rounds that have turned up 0 new
	// elements at the listing endpoint.
	blanks int
	// blankThreshold is the number of blanks a monitor will tolerate before
	// suspecting its tip is broken (e.g.post was deleted).
	blankThreshold int
	// tip is a slice of reddit thing names, the first of which represents
	// the "tip", which the monitor uses to requests new posts by using it
	// as a reference point (i.e.asks Reddit for posts "after" the tip).
	tip []string
	// path is the listing endpoint the monitor monitors. This path is
	// appended to the reddit monitor url (e.g./user/robert).
	path string

	replyHandler ReplyHandler
	postHandler  PostHandler
	userHandler  UserHandler
}

// monitorFromPath provides a monitor monitor from the listing endpoint.
func monitorFromPath(
	scrape Scraper,
	path string,
	handlePost postHandler,
	handleComment commentHandler,
	handleMessage messageHandler,
) (Monitor, error) {
	if handlePost == nil && handleComment == nil && handleMessage == nil {
		return nil, fmt.Errorf("no handlers provided for events")
	}

	b := &monitor{
		handlePost:    handlePost,
		handleComment: handleComment,
		handleMessage: handleMessage,
		path:          path,
		tip:           []string{""},
	}

	if err := m.sync(scrape); err != nil {
		return nil, err
	}

	return b, nil
}

// dispatch starts goroutines of the appropriate handler function for all new
// elements the monitor discovers.
func (m *monitor) dispatch(
	posts []*redditproto.Link,
	comments []*redditproto.Comment,
	messages []*redditproto.Message,
) {
	if m.handlePost != nil {
		for _, post := range posts {
			go m.handlePost(post)
		}
	}

	if m.handleComment != nil {
		for _, comment := range comments {
			go m.handleComment(comment)
		}
	}

	if m.handleMessage != nil {
		for _, message := range messages {
			go m.handleMessage(message)
		}
	}
}

// Update checks for new content at the monitored listing endpoint and forwards
// new content to the bot for processing.
func (m *monitor) Update(scrape Scraper, probe Prober) error {
	posts, comments, messages, err := scrape(m.path, m.tip[0], -1)
	if err != nil {
		return err
	}

	m.dispatch(posts, comments, messages)
	return m.updateTip(posts, comments, messages, probe)
}

// sync fetches the current tip of a listing endpoint, so that grawbots crawling
// forward in time don't treat it as a new post, or reprocess it when restarted.
func (m *monitor) sync(scrape Scraper) error {
	posts, messages, comments, err := scrape(m.path, "", 1)
	if err != nil {
		return err
	}

	things := merge(posts, messages, comments)
	if len(things) == 1 {
		m.tip = []string{things[0].GetName()}
	} else {
		m.tip = []string{""}
	}

	return nil
}

//.Tip.s the monitor's list of names from the endpoint listing it
// uses to keep track of its position in the monitored listing (e.g.a user's
// page or its position in a subreddit's history).
func (m *monitor) updateTip(
	posts []*redditproto.Link,
	comments []*redditproto.Comment,
	messages []*redditproto.Message,
	probe Prober,
) error {
	things := merge(posts, comments, messages)
	if len(things) == 0 {
		return m.healthCheck(probe)
	}

	names := make([]string, len(things))
	for i := 0; i < len(things); i++ {
		names[i] = things[i].GetName()
	}
	m.tip = append(names, m.tip...)
	if len(m.tip) > maxTipSize {
		m.tip = m.tip[0:maxTipSize]
	}

	return nil
}

// healthCheck checks the health of the tip when nothing is returned from a
// scrape enough times.
func (m *monitor) healthCheck(probe Prober) error {
	m.blanks++
	if m.blanks > m.blankThreshold {
		m.blanks = 0
		broken, err := m.fixTip(probe)
		if err != nil {
			return err
		}
		if !broken {
			m.blankThreshold += blankThreshold
		}
	}
	return nil
}

// fixTip checks that the fullname at the front of the tip is still valid (e.g.
// not deleted).If it isn't, it shaves the tip.fixTip returns whether the tip
// was broken.
func (m *monitor) fixTip(probe Prober) (bool, error) {
	exists, err := probe(m.tip[0])
	if err != nil {
		return false, err
	}

	if exists == false {
		m.shaveTip()
	}

	return !exists, nil
}

// shaveTip shaves the latest fullname off of the tip, promoting the preceding
// fullname if there is one or resetting the tip if there isn't.
func (m *monitor) shaveTip() {
	if len(m.tip) <= 1 {
		m.tip = []string{""}
	} else {
		m.tip = m.tip[1:]
	}
}
