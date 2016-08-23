package monitor

import (
	"strings"

	"github.com/turnage/redditproto"
)

// PostMonitor returns a monitor for new posts in a subreddit(s).
func PostMonitor(
	scrape Scraper,
	handlePost postHandler,
	subreddits []string,
) (Monitor, error) {
	return baseFromPath(
		scrape,
		"/r/"+strings.Join(subreddits, "+")+"/new",
		handlePost,
		nil,
		nil,
	)
}

// UserMonitor returns a monitor for new posts or comments by a user.
func UserMonitor(
	scrape Scraper,
	handlePost postHandler,
	handleComment commentHandler,
	user string,
) (Monitor, error) {
	return baseFromPath(
		scrape,
		"/user/"+user,
		handlePost,
		handleComment,
		nil,
	)
}

// CommentReplyMonitor returns a monitor for new replies to the bot's comments.
func CommentReplyMonitor(
	scrape Scraper,
	handleComment commentHandler,
) (Monitor, error) {
	return baseFromPath(
		scrape,
		"/message/comments",
		nil,
		handleComment,
		nil,
	)
}

// PostReplyMonitor returns a monitor for new replies to the bot's posts.
func PostReplyMonitor(
	scrape Scraper,
	handleComment commentHandler,
) (Monitor, error) {
	return baseFromPath(
		scrape,
		"/message/selfreply",
		nil,
		handleComment,
		nil,
	)
}

// MentionMonitor returns a monitor for new mentions of the bot's username
// across Reddit.
func MentionMonitor(
	scrape Scraper,
	handleComment commentHandler,
) (Monitor, error) {
	return baseFromPath(
		scrape,
		"/message/mentions",
		nil,
		handleComment,
		nil,
	)
}

type messageMonitor struct {
	base
	handleMessage messageHandler
}

// MessageMonitor returns a monitor for new private messages to the bot.
func MessageMonitor(
	scrape Scraper,
	handleMessage messageHandler,
) (Monitor, error) {
	mon := &messageMonitor{
		handleMessage: handleMessage,
	}

	b, err := baseFromPath(
		scrape,
		"/message/inbox",
		nil,
		nil,
		handleMessage,
	)

	if err == nil {
		mon.base = *(b.(*base))
	}

	return mon, err
}

// dispatchFilter only dispatches messages that were originally messages to the
// message handler.
func (m *messageMonitor) dispatchFilter(message *redditproto.Message) {
	if !message.GetWasComment() {
		m.handleMessage(message)
	}
}
