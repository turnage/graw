// Package monitor includes monitors for different parts of Reddit, such as a
// user inbox or a subreddit's post feed.
package monitor

import (
	"fmt"

	"github.com/turnage/graw/internal/monitor/internal/handlers"
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

// Direction represents a direction in time (Forward, or Backward).
type Direction int

const (
	Forward  = iota
	Backward = iota
)

const (
	// The blank threshold is the amount of.s returning 0 new
	// elements in the monitored listing the monitor will tolerate before
	// suspecting the tip of the listing has been deleted.
	blankThreshold = 2
	// maxTipSize is the maximum size of the tip log (number of backup tips
	// + the current tip).
	maxTipSize = 10
)

// Monitor defines the controls for a Monitor.
type Monitor interface {
	// Update will check for new events, and send them to the Monitor's
	// handler.
	Update(operator.Operator) error
}

// Monitors returns all of the monitors suitable for the bot. Some monitors
// can't be provided by this (such as user monitors) because the bot must
// request them.
func Monitors(
	bot interface{},
	subreddits []string,
	op operator.Operator,
	dir Direction,
) ([]Monitor, error) {
	monitors := []Monitor{}
	if han, ok := bot.(handlers.PostHandler); ok && len(subreddits) > 0 {
		mon, err := PostMonitor(op, han, subreddits, dir)
		if err != nil {
			return nil, err
		}
		monitors = append(monitors, mon)
	}

	messageHandler, _ := bot.(handlers.MessageHandler)
	postReplyHandler, _ := bot.(handlers.PostReplyHandler)
	commentReplyHandler, _ := bot.(handlers.CommentReplyHandler)
	mentionHandler, _ := bot.(handlers.MentionHandler)

	if messageHandler != nil ||
		postReplyHandler != nil ||
		commentReplyHandler != nil ||
		mentionHandler != nil {
		mon, err := InboxMonitor(
			op,
			messageHandler,
			postReplyHandler,
			commentReplyHandler,
			mentionHandler,
			dir,
		)
		if err != nil {
			return nil, err
		}
		monitors = append(monitors, mon)
	}

	return monitors, nil
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
	handlePost func(*redditproto.Link)
	// handleComment is the function the monitor uses to handle new comments
	// it finds.
	handleComment func(*redditproto.Comment)
	// handleMessage is the function the monitor uses to handle new messages
	// it finds.
	handleMessage func(*redditproto.Message)
	// dir is the direction in time the monitor monitors reddit.
	dir Direction
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
	// path is the listing endpoint the monitor monitors.This path is
	// appended to the reddit base url (e.g./user/robert).
	path string
}

// dispatch starts goroutines of the appropriate handler function for all new
// elements the monitor discovers.
func (b *base) dispatch(
	posts []*redditproto.Link,
	comments []*redditproto.Comment,
	messages []*redditproto.Message,
) error {
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

	return nil
}

// Update checks for new content at the monitored listing endpoint and forwards
// new content to the bot for processing.
func (b *base) Update(op operator.Operator) error {
	after := ""
	before := ""

	if b.dir == Forward {
		before = b.tip[0]
	} else if b.dir == Backward {
		after = b.tip[0]
	}

	posts, comments, messages, err := op.Scrape(
		b.path,
		after,
		before,
		operator.MaxLinks,
	)
	if err != nil {
		return err
	}

	b.dispatch(posts, comments, messages)
	return b.updateTip(posts, comments, messages, op)
}

// sync fetches the current tip of a listing endpoint, so that grawbots crawling
// forward in time don't treat it as a new post, or reprocess it when restarted.
func (b *base) sync(op operator.Operator) error {
	posts, messages, comments, err := op.Scrape(b.path, "", "", 1)
	if err != nil {
		return err
	}

	things := merge(posts, messages, comments, b.dir)
	if len(things) != 1 {
		return fmt.Errorf("invalid path")
	}

	b.tip = []string{things[0].GetName()}

	return nil
}

//.Tip.s the monitor's list of names from the endpoint listing it
// uses to keep track of its position in the monitored listing (e.g.a user's
// page or its position in a subreddit's history).
func (b *base) updateTip(
	posts []*redditproto.Link,
	comments []*redditproto.Comment,
	messages []*redditproto.Message,
	op operator.Operator,
) error {
	things := merge(posts, comments, messages, b.dir)
	if len(things) == 0 {
		return b.healthCheck(op)
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
func (b *base) healthCheck(op operator.Operator) error {
	b.blanks++
	if b.blanks > b.blankThreshold {
		b.blanks = 0
		broken, err := b.fixTip(op)
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
func (b *base) fixTip(op operator.Operator) (bool, error) {
	exists, err := op.IsThereThing(b.tip[0])
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
// one slice, ordered in creation time by dir.merge assumes all of the listings
// provided are ordered by timestamp according to dir.In general the total
// listing size will very rarely exceed 15 or so.
func merge(
	posts []*redditproto.Link,
	comments []*redditproto.Comment,
	messages []*redditproto.Message,
	dir Direction,
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
			if things[j].GetCreatedUtc() < things[j-1].GetCreatedUtc() {
				if dir == Forward {
					swap := things[j-1]
					things[j-1] = things[j]
					things[j] = swap
				}
			} else if dir == Backward {
				swap := things[j-1]
				things[j-1] = things[j]
				things[j] = swap
			}
		}
	}

	return things
}
