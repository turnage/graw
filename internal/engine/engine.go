// Package engine runs graw bots.
package engine

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/turnage/graw/internal/botfaces"
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

const (
	// blockTime is the amount of time to block between letting the next
	// monitor update.
	blockTime = time.Minute / 30
)

// Engine is the interface for the engine to the graw package.
type Engine interface {
	// Run runs the main engine loop, and returns an error if it encounters
	// one it can't handle. This runs indefinitely.
	Run() error
}

type base struct {
	// op is the operator the engine uses to make calls to Reddit.
	op operator.Operator
	// bot is the bot this engine runs.
	bot interface{}
	// dir is the direction in time this engine runs the bot.
	dir monitor.Direction
	// stopSig is a channel over which bots can send a signal to the engine
	// to stop.
	stopSig chan bool
	// stop is a flag for the engine to conclude its main loop.
	stop bool
	// Mutex protects all fields below.
	sync.Mutex
	// monitors is a list of the monitors this engine uses to get events
	// from Reddit.
	monitors *list.List
	// userMonitors is a map of username to the monitors dedicated to that
	// username.
	userMonitors map[string]*list.Element
}

// Reply submits a reply.
func (b *base) Reply(parentName, text string) error {
	return b.op.Reply(parentName, text)
}

// SendMessage sends a private message.
func (b *base) SendMessage(user, subject, text string) error {
	return b.op.Compose(user, subject, text)
}

// SelfPost makes a self (text) post to a subreddit.
func (b *base) SelfPost(subreddit, title, text string) error {
	return b.op.Submit(subreddit, "self", title, text)
}

// LinkPost makes a link post to a subreddit.
func (b *base) LinkPost(subreddit, title, url string) error {
	return b.op.Submit(subreddit, "link", title, url)
}

// WatchUser starts monitoring a useb.
func (b *base) WatchUser(user string) error {
	han, ok := b.bot.(botfaces.UserHandler)
	if !ok {
		return fmt.Errorf("bot cannot handle user posts or comments")
	}

	mon, err := monitor.UserMonitor(
		b.op,
		han.UserPost,
		han.UserComment,
		user,
		b.dir,
	)
	if err != nil {
		return err
	}

	b.Lock()
	defer b.Unlock()
	b.userMonitors[user] = b.monitors.PushBack(mon)
	return nil
}

// Unwatch users stops monitoring a useb.
func (b *base) UnwatchUser(user string) error {
	b.Lock()
	defer b.Unlock()

	if elem, ok := b.userMonitors[user]; ok {
		b.monitors.Remove(elem)
		delete(b.userMonitors, user)
	}

	return nil
}

// DigestThread returns a Link with a parsed comment tree.
func (b *base) DigestThread(permalink string) (*redditproto.Link, error) {
	return b.op.Thread(permalink)
}

// Stop stops the engine.
func (b *base) Stop() {
	b.stopSig <- true
}

func (b *base) Run() error {
	if actor, ok := b.bot.(botfaces.Actor); ok {
		actor.TakeEngine(b)
	}

	if loader, ok := b.bot.(botfaces.Loader); ok {
		loader.SetUp()
		defer loader.TearDown()
	}

	for !b.stop {
		select {
		case <-b.stopSig:
			b.stop = true
		case <-time.After(blockTime):
			if err := b.updateMonitors(); err != nil {
				if failer, ok := b.bot.(botfaces.Failer); !(ok && !failer.Fail(err)) {
					return err
				}
			}
		}
	}

	return nil
}

func (b *base) updateMonitors() error {
	b.Lock()
	monitors := []monitor.Monitor{}
	for i := b.monitors.Front(); i != nil; i = i.Next() {
		monitors = append(monitors, i.Value.(monitor.Monitor))
	}
	b.Unlock()

	for _, mon := range monitors {
		if err := mon.Update(b.op); err != nil {
			return err
		}
	}

	return nil
}
