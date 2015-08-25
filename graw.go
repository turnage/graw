// Package graw runs Reddit bots.
package graw

import (
	"strings"

	"github.com/turnage/graw/api"
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/operator"
)

// Run runs a bot against live reddit.
// agent should be the filename of a configured user agent protobuffer.
// graw will monitor all provided subreddits.
//
// For more information, see
// https://github.com/turnage/graw/wiki/Getting-Started
func Run(agent string, bot interface{}, subreddits ...string) error {
	op, err := operator.New(agent)
	if err != nil {
		return err
	}

	monitors := []monitor.Monitor{}
	if postHandler, ok := bot.(api.PostHandler); ok {
		monitors = append(
			monitors,
			&monitor.PostMonitor{
				Query: strings.Join(subreddits, "+"),
				Op:    op,
				Bot:   postHandler,
			},
		)
	}
	if inboxHandler, ok := bot.(api.InboxHandler); ok {
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

	actor, _ := bot.(api.Actor)
	loader, _ := bot.(api.Loader)
	failer, _ := bot.(api.Failer)
	return eng.Run(actor, loader, failer)
}
