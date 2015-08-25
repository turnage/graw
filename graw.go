package graw

import (
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/operator"
)

// Run runs a bot against live reddit. agent should be the filename of a
// configured user agent protobuffer. The bot will monitor all provide
// subreddits.
func Run(agent string, bot interface{}, subreddits ...string) error {
	op, err := operator.New(agent)
	if err != nil {
		return err
	}

	eng := &rtEngine{
		op:  op,
		mon: monitor.New(op, subreddits),
	}

	actor, _ := bot.(Actor)
	loader, _ := bot.(Loader)
	postHandler, _ := bot.(monitor.PostHandler)
	inboxHandler, _ := bot.(monitor.InboxHandler)
	return eng.Run(actor, loader, postHandler, inboxHandler)
}
