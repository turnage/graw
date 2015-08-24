package graw

import (
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

// rtEngine runs bots against real time Reddit.
type rtEngine struct {
	// bot is the bot this engine will run.
	bot Bot
	// op is the rtEngine's operator for making reddit api calls.
	op *operator.Operator
	// mon is the monitor rtEngine gets real time updates from.
	mon *monitor.Monitor

	// stop is a switch bots can set to signal the engine should stop.
	stop bool
}

// ReplyToMessage replies to an inboxable (private message, comment reply).
func (r *rtEngine) ReplyToInbox(msg *redditproto.Message, text string) error {
	return r.op.Reply(msg.GetName(), text)
}

// SendMessage sends a private message.
func (r *rtEngine) SendMessage(user, subject, text string) error {
	return r.op.Compose(user, subject, text)
}

// ReplyToPost posts a top-level comment on a submission.
func (r *rtEngine) ReplyToPost(post *redditproto.Link, text string) error {
	return r.op.Reply(post.GetName(), text)
}

// SelfPost makes a self (text) post to a subreddit.
func (r *rtEngine) SelfPost(subreddit, title, text string) error {
	return r.op.Submit(subreddit, "self", title, text)
}

// LinkPost makes a link post to a subreddit.
func (r *rtEngine) LinkPost(subreddit, title, url string) error {
	return r.op.Submit(subreddit, "link", title, url)
}

// ScrapeThread returns a full thread with prepared comment tree.
func (r *rtEngine) ScrapeThread(
	post *redditproto.Link,
) (*redditproto.Link, error) {
	return r.ScrapeThreadAt(post.GetPermalink())
}

// ScrapeThreadAt returns a full thread with prepared comment tree.
func (r *rtEngine) ScrapeThreadAt(permalink string) (*redditproto.Link, error) {
	return r.op.Thread(permalink)
}

// Stop is a function exposed to bots to stop the engine.
func (r *rtEngine) Stop() {
	r.stop = true
}

// Run is the main engine loop which runs the bot.
func (r *rtEngine) Run() error {
	r.bot.SetUp()
	defer r.bot.TearDown()

	go r.mon.Run()

	for !r.stop {
		select {
		case post := <-r.mon.NewPosts:
			go r.bot.Post(r, post)
		case message := <-r.mon.NewMessages:
			go r.bot.Message(r, message)
		case reply := <-r.mon.NewCommentReplies:
			go r.bot.Reply(r, reply)
		case mention := <-r.mon.NewMentions:
			go r.bot.Mention(r, mention)
		case err := <-r.mon.Errors:
			return err
		}
	}
	return nil
}
