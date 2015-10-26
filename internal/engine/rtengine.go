package engine

import (
	"fmt"

	"github.com/turnage/graw/api"
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

// rtEngine runs bots against real time Reddit.
type rtEngine struct {
	base
	// bot is the bot rtEngine is running.
	bot interface{}
	// monitors is a set of the monitors rtEngine gets events from.
	monitors map[string]monitor.Monitor
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
		base: base{
			op:     op,
			actor:  actor,
			failer: failer,
			loader: loader,
		},
		bot:      bot,
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
