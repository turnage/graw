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
// https://github.com/turnage/graw/wiki
func Run(agent string, bot interface{}, subreddits ...string) error {
	op, err := operator.New(agent)
	if err != nil {
		return err
	}

	actor, _ := bot.(api.Actor)
	failer, _ := bot.(api.Failer)
	loader, _ := bot.(api.Loader)
	eng := &rtEngine{
		Op:       op,
		Monitors: monitors(op, bot, subreddits),
		Actor:    actor,
		Failer:   failer,
		Loader:   loader,
	}

	return eng.Run()
}

// monitors returns the monitors appropriate for the given bot, based on the
// interfaces it implements.
func monitors(
	op operator.Operator,
	bot interface{},
	subreddits []string,
) []monitor.Monitor {
	mons := []monitor.Monitor{}
	if postHandler, ok := bot.(api.PostHandler); ok {
		mons = append(
			mons,
			&monitor.PostMonitor{
				Query: strings.Join(subreddits, "+"),
				Op:    op,
				Bot:   postHandler,
			},
		)
	}

	if mon := monitor.InboxMonitor(op, bot); mon != nil {
		mons = append(mons, mon)
	}
	return mons
}
