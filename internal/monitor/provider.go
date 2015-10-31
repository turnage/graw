package monitor

import (
	"strings"

	"github.com/turnage/graw/internal/monitor/internal/handlers"
	"github.com/turnage/graw/internal/operator"
)

// PostMonitor returns a monitor for new posts in a subreddit(s).
func PostMonitor(
	op operator.Operator,
	bot handlers.PostHandler,
	subreddits []string,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/r/"+strings.Join(subreddits, "+"),
		bot.Post,
		nil,
		nil,
		dir,
	)
}

// UserMonitor returns a monitor for new posts or comments by a user.
func UserMonitor(
	op operator.Operator,
	bot handlers.UserHandler,
	user string,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/user/"+user,
		bot.UserPost,
		bot.UserComment,
		nil,
		dir,
	)
}

// MessageMonitor returns a monitor for new private messages to the bot.
func MessageMonitor(
	op operator.Operator,
	bot handlers.MessageHandler,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/message/messages",
		nil,
		nil,
		bot.Message,
		dir,
	)
}

// CommentReplyMonitor returns a monitor for new replies to the bot's comments.
func CommentReplyMonitor(
	op operator.Operator,
	bot handlers.CommentReplyHandler,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/message/comments",
		nil,
		bot.CommentReply,
		nil,
		dir,
	)
}

// PostReplyMonitor returns a monitor for new replies to the bot's posts.
func PostReplyMonitor(
	op operator.Operator,
	bot handlers.PostReplyHandler,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/message/selfreply",
		nil,
		bot.PostReply,
		nil,
		dir,
	)
}

// MentionMonitor returns a monitor for new mentions of the bot's username
// across Reddit.
func MentionMonitor(
	op operator.Operator,
	bot handlers.MentionHandler,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/message/mentions",
		nil,
		bot.Mention,
		nil,
		dir,
	)
}
