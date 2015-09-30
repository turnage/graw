package monitor

import (
	"fmt"
	"strings"

	"github.com/turnage/graw/api"
	"github.com/turnage/graw/internal/monitor/internal/scanner"
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

// postMonitor monitors subreddits for new posts, and sends them to its handler.
type postMonitor struct {
	// postScanner is the scanner postMonitor uses to get updates from
	// subreddits it monitors.
	postScanner *scanner.Scanner
	// postHandler is the handler PostMonitor will send new posts to.
	postHandler api.PostHandler
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
		postScanner: scanner.New(
			fmt.Sprintf(
				"/r/%s/new",
				strings.Join(subreddits, "+"),
			),
			op,
		),
		postHandler: postHandler,
	}
}

// Update polls for new posts and sends them to Bot when they are found.
func (p *postMonitor) Update() error {
	postThings, err := p.postScanner.Scan()
	if err != nil {
		return err
	}

	posts := make([]*redditproto.Link, len(postThings))
	for i, thing := range postThings {
		posts[len(posts)-1-i] = thing.(*redditproto.Link)
	}

	for _, post := range posts {
		go p.postHandler.Post(post)
	}

	return nil
}
