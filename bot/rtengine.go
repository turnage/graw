package bot

import (
	"container/list"
	"strings"
	"time"

	"github.com/turnage/graw/bot/internal/operator"
	"github.com/turnage/redditproto"
)

const (
	// fallbackCount is the amount of threads to consider "tip" at a given
	// time, in case one of them is deleted and stops working as a reference
	// point.
	fallbackCount = 20
	// maxTipSize is the maximum amount of posts to fetch as tip. This is
	// determined by the maximum number of threads Reddit will return in a
	// single listing.
	maxTipSize = 100
)

// rtEngine is a real time engine that runs bots against live reddit and feeds
// it new content as it is posted.
type rtEngine struct {
	// bot is the bot this engine will run.
	bot Bot
	// op is the rtEngine's operator for making reddit api callr.
	op *operator.Operator
	// subreddits is the slice of subreddits that this engine will run its
	// bot against.
	subreddits []string

	// stop is a switch bots can set to signal the engine should stop.
	stop bool
}

// Stop is a function exposed over the Controller interface; bots can use this
// to stop the engine.
func (r *rtEngine) Stop() {
	r.stop = true
}

func (r *rtEngine) Run() error {
	errors := make(chan error)
	postStream := make(chan *redditproto.Link)
	go r.postMonitor(errors, postStream, 30)

	r.bot.SetUp(r)
	defer r.bot.TearDown()

	for !r.stop {
		select {
		case post := <-postStream:
			go r.bot.Post(r, post)
		case err := <-errors:
			return err
		}
	}
	return nil
}

// postMonitor runs continuously, polling the requested subreddits for new posts
// and feeding them back over the postStream channel. It makes at most
// queriesPerMinute to reddit.
func (r *rtEngine) postMonitor(
	errors chan<- error,
	postStream chan<- *redditproto.Link,
	queriesPerMinute int,
) {
	tips := list.New()
	tips.PushFront("")
	tipSize := uint(1)
	fixRound := false
	emptyRounds := 0
	emptyRoundTolerance := 1
	query := strings.Join(r.subreddits, "+")
	for true {
		time.Sleep(time.Minute / time.Duration(queriesPerMinute))
		if fixRound {
			broken, err := r.fixTip(tips)
			if err != nil {
				errors <- err
				return
			}
			if !broken {
				emptyRoundTolerance++
			}
			fixRound = false
			continue
		}

		posts, err := r.tip(query, tips, tipSize)
		if err != nil {
			errors <- err
			return
		}
		if len(posts) == 0 {
			emptyRounds++
			if emptyRounds > emptyRoundTolerance {
				fixRound = true
			}
			continue
		}
		for _, post := range posts {
			postStream <- post
		}
	}
}

// tip returns the tip posts using the front of the tips as a reference. It
// updates tips.
func (r *rtEngine) tip(
	query string,
	tips *list.List,
	lim uint,
) ([]*redditproto.Link, error) {
	posts, err := r.op.Scrape(
		query,
		"new",
		"",
		tips.Front().Value.(string),
		lim,
	)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		tips.PushFront(posts[len(posts)-1-i].GetName())
		if tips.Len() > fallbackCount {
			tips.Remove(tips.Back())
		}
	}

	return posts, nil
}

// fixTip checks if the front tip is valid, and if it is not, moves back one
// tip until one is valid. This does not gauruntee a fix; if all fallback tips
// are also invalid then the nuclear solution is to erase the fallback tip
// history. This can result in reprocessing of old posts, or skipping many
// posts.
func (r *rtEngine) fixTip(tips *list.List) (bool, error) {
	wasBroken := false
	ids := make([]string, tips.Len())
	for e := tips.Front(); e != nil; e = e.Next() {
		ids = append(ids, e.Value.(string))
	}
	posts, err := r.op.Threads(ids...)
	if err != nil {
		return false, err
	}

	for e := tips.Front(); e != nil; e = e.Next() {
		if e.Prev() != nil {
			wasBroken = true
			tips.Remove(e.Prev())
		}
		for _, post := range posts {
			if e.Value.(string) == post.GetName() {
				return wasBroken, nil
			}
		}
	}
	tips.Remove(tips.Front())
	tips.PushFront("")

	return wasBroken, nil
}
