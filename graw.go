package graw

import (
	"strings"

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

	monitors := []monitor.Monitor{}
	if postHandler, ok := bot.(monitor.PostHandler); ok {
		monitors = append(
			monitors,
			&monitor.PostMonitor{
				Query: strings.Join(subreddits, "+"),
				Op:    op,
				Bot:   postHandler,
			},
		)
	}
	if inboxHandler, ok := bot.(monitor.InboxHandler); ok {
		monitors = append(
			monitors,
			&monitor.InboxMonitor{
				Op:  op,
				Bot: inboxHandler,
			},
		)
	}

	eng := &rtEngine{
		op:       op,
		monitors: monitors,
	}

	actor, _ := bot.(Actor)
	loader, _ := bot.(Loader)
	failer, _ := bot.(Failer)
	return eng.Run(actor, loader, failer)
}
