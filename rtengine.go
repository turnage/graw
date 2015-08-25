package graw

import (
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

// rtEngine runs bots against real time Reddit.
type rtEngine struct {
	// op is the rtEngine's operator for making reddit api calls.
	op operator.Operator
	// monitors is a slice of the monitors rtEngine gets events from.
	monitors []monitor.Monitor

	// stop is a switch bots can set to signal the engine should stop.
	stop bool
}

// Reply submits a reply.
func (r *rtEngine) Reply(parentName, text string) error {
	return r.op.Reply(parentName, text)
}

// SendMessage sends a private message.
func (r *rtEngine) SendMessage(user, subject, text string) error {
	return r.op.Compose(user, subject, text)
}

// SelfPost makes a self (text) post to a subreddit.
func (r *rtEngine) SelfPost(subreddit, title, text string) error {
	return r.op.Submit(subreddit, "self", title, text)
}

// LinkPost makes a link post to a subreddit.
func (r *rtEngine) LinkPost(subreddit, title, url string) error {
	return r.op.Submit(subreddit, "link", title, url)
}

// DigestThread returns a Link with a parsed comment tree.
func (r *rtEngine) DigestThread(permalink string) (*redditproto.Link, error) {
	return r.op.Thread(permalink)
}

// Stop is a function exposed to bots to stop the engine.
func (r *rtEngine) Stop() {
	r.stop = true
}

// Run is the main engine loop.
func (r *rtEngine) Run(actor Actor, loader Loader, failer Failer) error {
	if loader != nil {
		loader.SetUp()
		defer loader.TearDown()
	}

	if actor != nil {
		actor.TakeEngine(r)
	}

	for !r.stop {
		for _, mon := range r.monitors {
			if err := mon.Update(); err != nil {
				if r.fail(failer, err) {
					return err
				}
			}
		}
	}

	return nil
}

// fail lets the bot decide whether to treat an error as a failure.
func (r *rtEngine) fail(failer Failer, err error) bool {
	if failer == nil {
		return false
	}

	return failer.Fail(err)
}
