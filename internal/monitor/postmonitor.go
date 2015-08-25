package monitor

import (
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

const (
	// maxTipSize is the number of posts to keep in the tracked tip. > 1
	// is kept because a tip is needed to fetch only posts newer than
	// that post. If one is deleted, PostMonitor moves to a fallback tip.
	maxTipSize = 15
)

// PostMonitor monitors subreddits for new posts, and sends them to its handler.
type PostMonitor struct {
	// Query is the multireddit query PostMonitor will use to find new posts
	// (e.g. self+funny).
	Query string
	// Posts is the number of posts PostMonitor has found since it began.
	Posts uint64
	// Bot is the handler PostMonitor will send new posts to.
	Bot PostHandler
	// Op is the operator through which the monitor will make update
	// requests to reddit.
	Op operator.Operator

	// tip is the list of latest posts in the monitored subreddits.
	tip []string
}

// Update polls for new posts and sends them to Bot when they are found.
func (p *PostMonitor) Update() error {
	if p.tip == nil {
		p.init()
	}

	posts, err := p.fetchTip()
	if err != nil {
		return err
	}

	if len(posts) == 0 {
		if err := p.fixTip(); err != nil {
			return err
		}
	}

	for _, post := range posts {
		go p.Bot.Post(post)
	}

	return nil
}

// init initializes the PostMonitor.
func (p *PostMonitor) init() {
	p.tip = make([]string, 1)
	p.tip[0] = ""
}

// fetchTip fetches the latest posts from the monitored subreddits. If there is
// no tip, fetchTip considers the call an adjustment round, and will fetch a new
// reference tip but discard the post (because, most likely, that post was
// already returned before).
func (p *PostMonitor) fetchTip() ([]*redditproto.Link, error) {
	tip := p.tip[len(p.tip)-1]
	links := uint(operator.MaxLinks)
	adjustment := false
	if tip == "" {
		links = 1
		adjustment = true
	}

	posts, err := p.Op.Scrape(
		p.Query,
		"new",
		"",
		tip,
		links,
	)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		p.tip = append(p.tip, posts[len(posts)-1-i].GetName())
	}

	if len(p.tip) > maxTipSize {
		p.tip = p.tip[len(p.tip)-maxTipSize:]
	}

	if adjustment && len(posts) == 1 {
		return nil, nil
	}

	return posts, nil
}

// fixTip attempts to fix the PostMonitor's reference point for new posts. If it
// has been deleted, fixTip will move to a fallback tip.
func (p *PostMonitor) fixTip() error {
	posts, err := p.Op.Threads(p.tip[len(p.tip)-1])
	if err != nil {
		return err
	}

	if len(posts) != 1 {
		p.shaveTip()
	}

	return nil
}

// shaveTip shaves off the latest tip thread name. If all tips are shaved off,
// uses an empty tip name (this will just get the latest threads).
func (p *PostMonitor) shaveTip() {
	if len(p.tip) == 1 {
		p.tip[0] = ""
		return
	}

	p.tip = p.tip[:len(p.tip)-1]
}
