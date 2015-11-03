// Package graw runs Reddit bots.
package graw

import (
	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/engine"
	"github.com/turnage/graw/internal/operator"
)

// Engine is the interface bots can use to request things from the Engine, like
// data from Reddit, that it makes a post, new event types, etc.
type Engine api.Engine

// Run runs a bot against live Reddit.
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

	eng, err := engine.RealTime(bot, op, subreddits)
	if err != nil {
		return err
	}

	return eng.Run()
}

// Scrape runs a bot against Reddit that has already happened, and moves
// backward in time.
// agent should be the filename of a configured user agent protobuffer.
// graw will scrape all provided subreddits.
//
// For more information, see
// https://github.com/turnage/graw/wiki
func Scrape(agent string, bot interface{}, subreddits ...string) error {
	op, err := operator.New(agent)
	if err != nil {
		return err
	}

	eng, err := engine.BackTime(bot, op, subreddits)
	if err != nil {
		return err
	}

	return eng.Run()
}
