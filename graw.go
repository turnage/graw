// Package graw runs Reddit bots.
package graw

import (
	"sync"

	"github.com/turnage/graw/internal/engine"
	"github.com/turnage/graw/internal/operator"
)

var (
	// mu protects engines.
	mu sync.Mutex
	// engines is a map of <k,v>:<bot,engine>.
	engines = map[interface{}]*engine.Engine{}
)

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

	return runEngine(bot, eng)
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

	return runEngine(bot, eng)
}

// GetEngine returns the engine running the given bot. The bot is used as a key
// to look up the corresponding Engine. If there is no engine for the bot, nil
// is returned.
func GetEngine(bot interface{}) Engine {
	mu.Lock()
	defer mu.Unlock()
	eng, ok := engines[bot]
	if ok {
		return eng
	}

	return nil
}

// runEngine runs an engine and manages its entry in the engine map.
func runEngine(bot interface{}, eng *engine.Engine) error {
	mu.Lock()
	engines[bot] = eng
	mu.Unlock()

	err := eng.Run()

	mu.Lock()
	delete(engines, bot)
	mu.Unlock()

	return err
}
