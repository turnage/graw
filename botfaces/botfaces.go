// Package botfaces defines interfaces graw uses to connect bots to event
// streams on Reddit. There is no need to import this package.
package botfaces

import (
	"github.com/turnage/graw/reddit"
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

// PostHandler defines methods for bots that handle new posts in
// subreddits they monitor.
type PostHandler interface {
	// Post is called when a post is made in a monitored subreddit that the
	// bot has not seen yet. [Called as goroutine.]
	Post(post *reddit.Post) error
}

// CommentHandler defines methods for bots that handle new comments in
// subreddits they monitor.
type CommentHandler interface {
	// Comment is called when a comment is made in a monitored subreddit
	// that the bot has not seen yet. [Called as goroutine.]
	Comment(post *reddit.Comment) error
}

// MessageHandler defines methods for bots that handle new private messages to
// their inbox.
type MessageHandler interface {
	// Message is called when the bot receives a new private message to its
	// account. [Called as goroutine.]
	Message(msg *reddit.Message) error
}

// PostReplyHandler defines methods for bots that handle new top-level comments
// on their submissions.
type PostReplyHandler interface {
	// Reply is called when the bot receives a reply to one of its
	// submissions in its inbox. [Called as goroutine.]
	//
	// The reply is in the form of a message because that is how it arrives
	// in the Reddit inbox, but it is originally a comment and replying to
	// this with the Reddit package will still generate a comment.
	PostReply(reply *reddit.Message) error
}

// CommentReplyHandler defines methods for bots that handle new comments made in
// reply to their own.
type CommentReplyHandler interface {
	// CommentReply is called when the bot receives a reply to one of its
	// comments in it its inbox. [Called as goroutine.]
	//
	// The reply is in the form of a message because that is how it arrives
	// in the Reddit inbox, but it is originally a comment and replying to
	// this with the Reddit package will still generate a comment.
	CommentReply(reply *reddit.Message) error
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
	//
	// The reply is in the form of a message because that is how it arrives
	// in the Reddit inbox, but it is originally a comment and replying to
	// this with the Reddit package will still generate a comment.
	Mention(mention *reddit.Message) error
}

// UserHandler defines methods for bots that handle activity by monitored users.
type UserHandler interface {
	// UserPost is called when a monitored user makes a post in a subreddit
	// the bot can view. [Called as goroutine.]
	UserPost(post *reddit.Post) error
	// UserComment is called when the monitored user makes a comment in a
	// subreddit the bot can view. [Called as goroutine.]
	UserComment(comment *reddit.Comment) error
}
