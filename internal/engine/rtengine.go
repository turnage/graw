package engine

import (
	"fmt"
	"sync"

	"github.com/turnage/graw/api"
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

// rtEngine runs bots against real time Reddit.
type rtEngine struct {
	// op is the rtEngine's operator for making reddit api calls.
	op operator.Operator
	// bot is the bot rtEngine is running.
	bot interface{}
	// actor is the bot's interface for receiving an interface to the
	// Engine, so that it can act through its Reddit account.
	actor api.Actor
	// failer is the bot's interface for handling errors. rtEngine will
	// defer to this to decide what to when it encounters an error.
	failer api.Failer
	// loader is the bot's interface for setting up and tearing down
	// resources.
	loader api.Loader
	// mu protects all variable below.
	mu sync.Mutex
	// monitors is a set of the monitors rtEngine gets events from.
	monitors map[string]monitor.Monitor
	// stop is a switch bots can set to signal the engine should stop.
	stop bool
}

// New returns the ignition to a real time engine, so that it can be started.
// A real time engine runs a bot against Reddit as it happens, indefinitely,
// feeding it new events as the occur and allowing the bot to interact by making
// posts, sending messages, and more.
func RealTime(
	agent string,
	bot interface{},
	subreddits []string,
) (
	Ignition,
	error,
) {
	if bot == nil {
		return nil, fmt.Errorf("bot was nil")
	}

	op, err := operator.New(agent)
	if err != nil {
		return nil, err
	}

	actor, _ := bot.(api.Actor)
	failer, _ := bot.(api.Failer)
	loader, _ := bot.(api.Loader)
	monitors := make(map[string]monitor.Monitor)

	if postHandler, ok := bot.(api.PostHandler); ok && len(subreddits) > 0 {
		monitors["/r/"] = monitor.PostMonitor(
			op,
			postHandler,
			subreddits,
		)
	}

	messageHandler, _ := bot.(api.MessageHandler)
	postReplyHandler, _ := bot.(api.PostReplyHandler)
	commentReplyHandler, _ := bot.(api.CommentReplyHandler)
	mentionHandler, _ := bot.(api.MentionHandler)

	if messageHandler != nil ||
		postReplyHandler != nil ||
		commentReplyHandler != nil ||
		mentionHandler != nil {
		monitors["/messages/"] = monitor.InboxMonitor(
			op,
			messageHandler,
			postReplyHandler,
			commentReplyHandler,
			mentionHandler,
		)
	}

	return &rtEngine{
		op:       op,
		bot:      bot,
		actor:    actor,
		failer:   failer,
		loader:   loader,
		monitors: monitors,
	}, nil
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

// WatchUser starts monitoring a user.
func (r *rtEngine) WatchUser(user string) error {
	r.mu.Lock()
	if userHandler, ok := r.bot.(api.UserHandler); ok {
		r.monitors[user] = monitor.UserMonitor(r.op, userHandler, user)
	}
	r.mu.Unlock()
	return nil
}

// Unwatch users stops monitoring a user.
func (r *rtEngine) UnwatchUser(user string) error {
	r.mu.Lock()
	delete(r.monitors, user)
	r.mu.Unlock()
	return nil
}

// DigestThread returns a Link with a parsed comment tree.
func (r *rtEngine) DigestThread(permalink string) (*redditproto.Link, error) {
	return r.op.Thread(permalink)
}

// Stop stops the engine.
func (r *rtEngine) Stop() {
	r.mu.Lock()
	r.stop = true
	r.mu.Unlock()
}

// Run is the main engine loop.
func (r *rtEngine) Run() error {
	r.setup()
	for !r.stop {
		for _, mon := range r.monitors {
			if err := mon.Update(); err != nil {
				if r.fail(err) {
					return err
				}
			}
		}
	}

	return nil
}

// setup prepares the engine and bot to run.
func (r *rtEngine) setup() {
	if r.loader != nil {
		r.loader.SetUp()
		defer r.loader.TearDown()
	}

	if r.actor != nil {
		r.actor.TakeEngine(r)
	}
}

// fail lets the bot decide whether to treat an error as a failure.
func (r *rtEngine) fail(err error) bool {
	if r.failer == nil {
		return true
	}

	return r.failer.Fail(err)
}
