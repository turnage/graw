package monitor

import (
	"strings"

	"github.com/turnage/graw/internal/monitor/internal/handlers"
	"github.com/turnage/graw/internal/operator"
)

// postMonitor monitors subreddits for new posts, and sends them to its handler.
type postMonitor struct {
	base
}

// PostMonitor returns a post monitor for the requested subreddits, using bot
// to handle new posts it finds.
func PostMonitor(
	op operator.Operator,
	bot handlers.PostHandler,
	subreddits []string,
	dir Direction,
) (Monitor, error) {
	p := &postMonitor{
		base: base{
			handlePost: bot.Post,
			dir:        dir,
			path:       "/r/" + strings.Join(subreddits, "+"),
			tip:        []string{""},
		},
	}

	if dir == Forward {
		if err := p.sync(op); err != nil {
			return nil, err
		}
	}

	return p, nil
}
