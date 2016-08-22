// Package graw runs Reddit bots.
package graw

import (
	"sync"

	"github.com/turnage/graw/internal/client"
	"github.com/turnage/graw/internal/engine"
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
	cli, err := client.New(agent)
	if err != nil {
		return err
	}

	eng, err := engine.New(bot, cli, subreddits)
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
