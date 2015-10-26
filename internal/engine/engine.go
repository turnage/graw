// Package engine provides implementations for bot engines. See the provider
// functions for details about what context they run a bot in.
package engine

import (
	"sync"

	"github.com/turnage/graw/api"
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

// Ignition provides an interface to start the engine.
type Ignition interface {
	// Run should be called once to start the engine. It may run forever.
	Run() error
}

// base contains the base fields all engines use.
type base struct {
	// op is the operator engines use to make api calls to Reddit.
	op operator.Operator
	// actor is the bot's interface for receiving an interface to the
	// Engine, so that it can act through its Reddit account.
	actor api.Actor
	// failer is the bot's interface for handling errors. rtEngine will
	// defer to this to decide what to when it encounters an errob.
	failer api.Failer
	// loader is the bot's interface for setting up and tearing down
	// resources.
	loader api.Loader
	// mu protects all variables below.
	mu sync.Mutex
	// stop is a switch bots can set to signal the engine should stop.
	stop bool
}

// Reply is a noop.
func (b *base) Reply(parentName, text string) error {
	return nil
}

// SendMessage is a noop.
func (b *base) SendMessage(user, subject, text string) error {
	return nil
}

// SelfPost is a noop.
func (b *base) SelfPost(subreddit, title, text string) error {
	return nil
}

// LinkPost is a noop.
func (b *base) LinkPost(subreddit, title, url string) error {
	return nil
}

// WatchUser is a noop.
func (b *base) WatchUser(user string) error {
	return nil
}

// Unwatch is a noop.
func (b *base) UnwatchUser(user string) error {
	return nil
}

// DigestThread is a noop.
func (b *base) DigestThread(permalink string) (*redditproto.Link, error) {
	return nil, nil
}

// Stop stops the engine.
func (b *base) Stop() {
	b.mu.Lock()
	b.stop = true
	b.mu.Unlock()
}

// setup prepares the engine and bot to run.
func (b *base) setup() {
	if b.loader != nil {
		b.loader.SetUp()
		defer b.loader.TearDown()
	}

	if b.actor != nil {
		b.actor.TakeEngine(b)
	}
}

// fail lets the bot decide whether to treat an error as a failure.
func (b *base) fail(err error) bool {
	if b.failer == nil {
		return true
	}

	return b.failer.Fail(err)
}
