package graw

import (
	"log"
)

// Config configures a graw run or scan by specifying event sources. Each event
// type has a corresponding handler defined in graw/botfaces. The bot must be
// able to handle requested event types.
type Config struct {
	// New posts in all subreddits named here will be forwarded to the bot's
	// PostHandler.
	Subreddits []string
	// New comments in all subreddits named here will be forwarded to the
	// bot's CommentHandler.
	SubredditComments []string
	// New posts and comments made by all users named here will be forwarded
	// to the bot's UserHandler. Note that since a separate monitor must be
	// construced for every user, unlike subreddits, subscribing to the
	// actions of many users can delay updates from other event sources.
	Users []string
	// When true, replies to posts made by the bot's account will be
	// forwarded to the bot's PostReplyHandler.
	PostReplies bool
	// When true, replies to comments made by the bot's account will be
	// forwarded to the bot's CommentReplyHandler.
	CommentReplies bool
	// When true, mentions of the bot's username  will be forwarded to the
	// bot's MentionHandler.
	Mentions bool
	// When true, messages sent to the bot's inbox will be forwarded to the
	// bot's MessageHandler.
	Messages bool
	// If set, internal messages will be logged here. This is a spammy log
	// used for debugging graw.
	Logger *log.Logger
}
