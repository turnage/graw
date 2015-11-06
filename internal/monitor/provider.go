package monitor

import (
	"strings"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

// PostMonitor returns a monitor for new posts in a subreddit(s).
func PostMonitor(
	op operator.Operator,
	handlePost postHandler,
	subreddits []string,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/r/"+strings.Join(subreddits, "+")+"/new",
		handlePost,
		nil,
		nil,
		dir,
	)
}

// UserMonitor returns a monitor for new posts or comments by a user.
func UserMonitor(
	op operator.Operator,
	handlePost postHandler,
	handleComment commentHandler,
	user string,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/user/"+user,
		handlePost,
		handleComment,
		nil,
		dir,
	)
}

// CommentReplyMonitor returns a monitor for new replies to the bot's comments.
func CommentReplyMonitor(
	op operator.Operator,
	handleComment commentHandler,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/message/comments",
		nil,
		handleComment,
		nil,
		dir,
	)
}

// PostReplyMonitor returns a monitor for new replies to the bot's posts.
func PostReplyMonitor(
	op operator.Operator,
	handleComment commentHandler,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/message/selfreply",
		nil,
		handleComment,
		nil,
		dir,
	)
}

// MentionMonitor returns a monitor for new mentions of the bot's username
// across Reddit.
func MentionMonitor(
	op operator.Operator,
	handleComment commentHandler,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/message/mentions",
		nil,
		handleComment,
		nil,
		dir,
	)
}

type messageMonitor struct {
	base
	handleMessage messageHandler
}

// MessageMonitor returns a monitor for new private messages to the bot.
func MessageMonitor(
	op operator.Operator,
	handleMessage messageHandler,
	dir Direction,
) (Monitor, error) {
	mon := &messageMonitor{
		handleMessage: handleMessage,
	}

	b, err := baseFromPath(
		op,
		"/message/inbox",
		nil,
		nil,
		handleMessage,
		dir,
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
