package monitor

import (
	"strings"

	"github.com/turnage/graw/api"
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

const (
	// maxTipSize is the number of posts to keep in the tracked tip. > 1
	// is kept because a tip is needed to fetch only posts newer than
	// that post. If one is deleted, PostMonitor moves to a fallback tip.
	maxTipSize = 15
)

// postMonitor monitors subreddits for new posts, and sends them to its handler.
type postMonitor struct {
	// op is the operator through which the monitor will make update
	// requests to reddit.
	op operator.Operator
	// postHandler is the handler PostMonitor will send new posts to.
	postHandler api.PostHandler
	// query is the multireddit query PostMonitor will use to find new posts
	// (e.g. self+funny).
	query string
	// tip is the list of latest posts in the monitored subreddits.
	tip []string
}

// PostMonitor returns a post monitor for the requested subreddits, busing bot
// to handle new posts it finds. If bot cannot handle posts or there are no
// subreddits to monitor, returns nil.
func PostMonitor(
	op operator.Operator,
	bot interface{},
	subreddits []string,
) Monitor {
	postHandler, ok := bot.(api.PostHandler)
	if !ok {
		return nil
	}

	if len(subreddits) == 0 {
		return nil
	}

	return &postMonitor{
		op:          op,
		postHandler: postHandler,
		tip:         []string{""},
		query:       strings.Join(subreddits, "+"),
	}
}

// Update polls for new posts and sends them to Bot when they are found.
func (p *postMonitor) Update() error {
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
		go p.postHandler.Post(post)
	}

	return nil
}

// fetchTip fetches the latest posts from the monitored subreddits. If there is
// no tip, fetchTip considers the call an adjustment round, and will fetch a new
// reference tip but discard the post (because, most likely, that post was
// already returned before).
func (p *postMonitor) fetchTip() ([]*redditproto.Link, error) {
	tip := p.tip[len(p.tip)-1]
	links := uint(operator.MaxLinks)
	adjustment := false
	if tip == "" {
		links = 1
		adjustment = true
	}

	posts, err := p.op.Scrape(
		p.query,
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
func (p *postMonitor) fixTip() error {
	posts, err := p.op.Threads(p.tip[len(p.tip)-1])
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
func (p *postMonitor) shaveTip() {
	if len(p.tip) == 1 {
		p.tip[0] = ""
		return
	}

	p.tip = p.tip[:len(p.tip)-1]
}
