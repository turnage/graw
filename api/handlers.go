package api

import (
	"github.com/turnage/redditproto"
)

// It is recommended, though not required, to implement all inbox handlers if
// one is implemented. All unread inbox items are fetched in updates (Reddit
// does not offer filtering by type (mentions, post replies, etc)).
//
// Ex. If Mentions are handled, but not Messages, any unread messages will be
// wasting your network data. graw will not mark an unhandled inbox item as
// read.
//
// It shouldn't be a concern if your bot cannot receive the events it does not
// handle. If it does not submit, it cannot receive post replies. If it does not
// comment, it cannot receive comment replies.

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
	PostReply(reply *redditproto.Message)
}

// CommentReplyHandler defines methods for bots that handle new comments made in
// reply to their own.
type CommentReplyHandler interface {
	// CommentReply is called when the bot receives a reply to one of its
	// comments in it its inbox. [Called as goroutine.]
	CommentReply(reply *redditproto.Message)
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
	Mention(mention *redditproto.Message)
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
