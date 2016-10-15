// Package monitor includes monitors for different parts of Reddit, such as a
// user inbox or a subreddit's post feed.
package monitor

import (
	"fmt"

	"github.com/turnage/redditproto"
)

type commentHandler func(*redditproto.Comment)
type postHandler func(*redditproto.Link)
type messageHandler func(*redditproto.Message)

const (
	// The blank threshold is the amount of updates returning 0 new
	// elements in the monitored listing the monitor will tolerate before
	// suspecting the tip of the listing has been deleted.
	blankThreshold = 2
	// maxTipSize is the maximum size of the tip log (number of backup tips
	// + the current tip).
	maxTipSize = 10
)

// Scraper defines a function that takes a Reddit listing path, and a an id of
// an element within the listing, and returns all the elements following that
// element. If limit is <= 0, as many elements as possible should be returned.
type Scraper func(
	path,
	tip string,
	limit int,
) (
	[]*redditproto.Link,
	[]*redditproto.Comment,
	[]*redditproto.Message,
	error,
)

// Prober defines a function type that takes the id of a reddit thing and
// returns whether it exists.
type Prober func(id string) (bool, error)

// Monitor defines the controls for a Monitor.
type Monitor interface {
	// Update will check for new events, and send them to the Monitor's
	// handlers.
	Update(scrape Scraper, probe Prober) error
}

// redditThing is an interface for accessing attributes of Reddit types that
// implement the Redddit "Thing" class.
type redditThing interface {
	// GetName returns the fullname of the thing.
	GetName() string
	// GetCreatedUtc returns the creation time (Unix time) of the thing.
	GetCreatedUtc() float64
}

// base describes the core of any monitor; most monitors will use its fields and
// methods.
type base struct {
	// handlePost is the function the monitor uses to handle new posts it
	// finds.
	handlePost postHandler
	// handleComment is the function the monitor uses to handle new comments
	// it finds.
	handleComment commentHandler
	// handleMessage is the function the monitor uses to handle new messages
	// it finds.
	handleMessage messageHandler
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
	// appended to the reddit base url (e.g./user/robert).
	path string
}

// baseFromPath provides a monitor base from the listing endpoint.
func baseFromPath(
	scrape Scraper,
	path string,
	handlePost postHandler,
	handleComment commentHandler,
	handleMessage messageHandler,
) (Monitor, error) {
	if handlePost == nil && handleComment == nil && handleMessage == nil {
		return nil, fmt.Errorf("no handlers provided for events")
	}

	b := &base{
		handlePost:    handlePost,
		handleComment: handleComment,
		handleMessage: handleMessage,
		path:          path,
		tip:           []string{""},
	}

	if err := b.sync(scrape); err != nil {
		return nil, err
	}

	return b, nil
}

// dispatch starts goroutines of the appropriate handler function for all new
// elements the monitor discovers.
func (b *base) dispatch(
	posts []*redditproto.Link,
	comments []*redditproto.Comment,
	messages []*redditproto.Message,
) {
	if b.handlePost != nil {
		for _, post := range posts {
			go b.handlePost(post)
		}
	}

	if b.handleComment != nil {
		for _, comment := range comments {
			go b.handleComment(comment)
		}
	}

	if b.handleMessage != nil {
		for _, message := range messages {
			go b.handleMessage(message)
		}
	}
}

// Update checks for new content at the monitored listing endpoint and forwards
// new content to the bot for processing.
func (b *base) Update(scrape Scraper, probe Prober) error {
	posts, comments, messages, err := scrape(b.path, b.tip[0], -1)
	if err != nil {
		return err
	}

	b.dispatch(posts, comments, messages)
	return b.updateTip(posts, comments, messages, probe)
}

// sync fetches the current tip of a listing endpoint, so that grawbots crawling
// forward in time don't treat it as a new post, or reprocess it when restarted.
func (b *base) sync(scrape Scraper) error {
	posts, messages, comments, err := scrape(b.path, "", 1)
	if err != nil {
		return err
	}

	things := merge(posts, messages, comments)
	if len(things) == 1 {
		b.tip = []string{things[0].GetName()}
	} else {
		b.tip = []string{""}
	}

	return nil
}

//.Tip.s the monitor's list of names from the endpoint listing it
// uses to keep track of its position in the monitored listing (e.g.a user's
// page or its position in a subreddit's history).
func (b *base) updateTip(
	posts []*redditproto.Link,
	comments []*redditproto.Comment,
	messages []*redditproto.Message,
	probe Prober,
) error {
	things := merge(posts, comments, messages)
	if len(things) == 0 {
		return b.healthCheck(probe)
	}

	names := make([]string, len(things))
	for i := 0; i < len(things); i++ {
		names[i] = things[i].GetName()
	}
	b.tip = append(names, b.tip...)
	if len(b.tip) > maxTipSize {
		b.tip = b.tip[0:maxTipSize]
	}

	return nil
}

// healthCheck checks the health of the tip when nothing is returned from a
// scrape enough times.
func (b *base) healthCheck(probe Prober) error {
	b.blanks++
	if b.blanks > b.blankThreshold {
		b.blanks = 0
		broken, err := b.fixTip(probe)
		if err != nil {
			return err
		}
		if !broken {
			b.blankThreshold += blankThreshold
		}
	}
	return nil
}

// fixTip checks that the fullname at the front of the tip is still valid (e.g.
// not deleted).If it isn't, it shaves the tip.fixTip returns whether the tip
// was broken.
func (b *base) fixTip(probe Prober) (bool, error) {
	exists, err := probe(b.tip[0])
	if err != nil {
		return false, err
	}

	if exists == false {
		b.shaveTip()
	}

	return !exists, nil
}

// shaveTip shaves the latest fullname off of the tip, promoting the preceding
// fullname if there is one or resetting the tip if there isn't.
func (b *base) shaveTip() {
	if len(b.tip) <= 1 {
		b.tip = []string{""}
	} else {
		b.tip = b.tip[1:]
	}
}

// merge merges elements of multiple listings which implement redditThing into
// one slice, ordered in creation time. merge assumes all of the listings
// provided are ordered independently by creation time. The total listing size
// will very rarely exceed 15 or so.
func merge(
	posts []*redditproto.Link,
	comments []*redditproto.Comment,
	messages []*redditproto.Message,
) []redditThing {
	// Why is handling interface slices so painful?
	things := []redditThing{}
	for _, post := range posts {
		things = append(things, post)
	}
	for _, comment := range comments {
		things = append(things, comment)
	}
	for _, message := range messages {
		things = append(things, message)
	}

	for i := 0; i < len(things); i++ {
		for j := len(things) - 1; j > i; j-- {
			if things[j].GetCreatedUtc() > things[j-1].GetCreatedUtc() {
				swap := things[j-1]
				things[j-1] = things[j]
				things[j] = swap
			}
		}
	}

	return things
}
