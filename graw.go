// Package graw runs Reddit bots.
package graw

import (
	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/engine"
)

// Engine is the interface bots can use to request things from the Engine, like
// data from Reddit, that it makes a post, new event types, etc.
type Engine api.Engine

// Run runs a bot against live reddit.
// agent should be the filename of a configured user agent protobuffer.
// graw will monitor all provided subreddits.
//
// For more information, see
// https://github.com/turnage/graw/wiki
func Run(agent string, bot interface{}, subreddits ...string) error {
	eng, err := engine.RealTime(agent, bot, subreddits)
	if err != nil {
		return err
	}

	return eng.Run()
}
