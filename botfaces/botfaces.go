// Package botfaces defines the interfaces a bot can have, visibile to the graw
// engine. These interfaces allow graw to infer what events the bot cares about
// and can handle.
package botfaces

import (
	"github.com/turnage/redditproto"
	"time"
)

// Loader defines methods for bots that use external resources or need to do
// initialization.
type Loader interface {
	// SetUp is the first method ever called on the bot, and it will be
	// allowed to finish before other methods are called. Bots should
	// load resources here. If an error is returned, the engine will not
	// start, and the error will propagate up.
	SetUp() error
}

// Tearer defines methods for bots that need to tear things down after their run
// is finished.
type Tearer interface {
	// TearDown is the last method ever called on the bot, and all other
	// method calls will finish before this method is called. Bots should
	// unload resources here.
	TearDown()
}

// Failer defines methods bots can use to control how the Engine responds to
// failures.
type Failer interface {
	// Fail will be called when the engine encounters an error. The bot can
	// return true to instruct the engine to fail, or false to instruct the
	// engine to try again.
	//
	// This method will be called in the main engine loop; the bot may
	// choose to pause here or do other things to respond to the failure
	// (e.g. pause for three hours to respond to Reddit down time).
	Fail(err error) bool
}

// PostHandler defines methods for bots that handle new posts in
// subreddits they monitor.
type PostHandler interface {
	// Post is called when a post is made in a monitored subreddit that the
	// bot has not seen yet. [Called as goroutine.]
	Post(post *redditproto.Link)
}

// MessageHandler defines methods for bots that handle new private messages to
// their inbox.
type MessageHandler interface {
	// Message is called when the bot receives a new private message to its
	// account. [Called as goroutine.]
	Message(msg *redditproto.Message)
}

// PostReplyHandler defines methods for bots that handle new top-level comments
// on their submissions.
type PostReplyHandler interface {
	// Reply is called when the bot receives a reply to one of its
	// submissions in its inbox. [Called as goroutine.]
	PostReply(reply *redditproto.Comment)
}

// CommentReplyHandler defines methods for bots that handle new comments made in
// reply to their own.
type CommentReplyHandler interface {
	// CommentReply is called when the bot receives a reply to one of its
	// comments in it its inbox. [Called as goroutine.]
	CommentReply(reply *redditproto.Comment)
}

// MentionHandler defines methods for bots that handle username mentions. These
// will only appear in the inbox if
//
//  1. The bot's account preferences are set to monitor username mentions.
//  2. The event type is not shadowed.
//
// Reddit will shadow a username mention under two other event types. If the
// bot's username is mentioned in a comment posted on its submission, the inbox
// item will be a post reply, not a mention. If the bot's username is mentioned
// in a comment posted in reply to one of the bot's comments, the inbox item
// will be a comment reply, not a mention. In both cases, this handler is
// unused.
type MentionHandler interface {
	// Mention is called when the bot receives a username mention in its
	// inbox. [Called as goroutine.]
	Mention(mention *redditproto.Comment)
}

// UserHandler defines methods for bots that handle activity by monitored users.
type UserHandler interface {
	// UserPost is called when a monitored user makes a post in a subreddit
	// the bot can view. [Called as goroutine.]
	UserPost(post *redditproto.Link)
	// UserComment is called when the monitored user makes a comment in a
	// subreddit the bot can view. [Called as goroutine.]
	UserComment(comment *redditproto.Comment)
}

// BlockTimer defines methods that allow bots to specify the interval to block
// between monitor updates.
type BlockTimer interface {
	// To comply with the reddit API rules, this should be >= 1*time.Second.
	// Default: 2*time.Second
	BlockTime() time.Duration
}
