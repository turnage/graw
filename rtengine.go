package graw

import (
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

// rtEngine runs bots against real time Reddit.
type rtEngine struct {
	// op is the rtEngine's operator for making reddit api calls.
	op *operator.Operator
	// mon is the monitor rtEngine gets real time updates from.
	mon *monitor.Monitor

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

// Run is the main engine loop which runs the bot.
func (r *rtEngine) Run(
	loader Loader,
	postHandler PostHandler,
	inboxHandler InboxHandler,
) error {
	if loader != nil {
		loader.SetUp()
		defer loader.TearDown()
	}

	go r.mon.Run()

	for !r.stop {
		select {
		case post := <-r.mon.NewPosts:
			go postHandler.Post(r, post)
		case message := <-r.mon.NewMessages:
			go inboxHandler.Message(r, message)
		case reply := <-r.mon.NewCommentReplies:
			go inboxHandler.CommentReply(r, reply)
		case reply := <-r.mon.NewPostReplies:
			go inboxHandler.PostReply(r, reply)
		case mention := <-r.mon.NewMentions:
			go inboxHandler.Mention(r, mention)
		case err := <-r.mon.Errors:
			return err
		}
	}
	return nil
}
