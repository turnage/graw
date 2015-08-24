package graw

import (
	"github.com/turnage/redditproto"
)

// Loader defines methods for bots that use external resources or need to do
// initialization.
type Loader interface {
	// SetUp is the first method ever called on the bot, and it will be
	// allowed to finish before other methods are called. Bots should
	// load resources here.
	SetUp() error
	// TearDown is the last method ever called on the bot, and all other
	// method calls will finish before this method is called. Bots should
	// unload resources here.
	TearDown() error
}

// PostHandler defines methods for bots that handle new posts in
// subreddits they monitor.
type PostHandler interface {
	// Post is called when a post is made in a monitored subreddit that the
	// bot has not seen yet. [Called as goroutine.]
	Post(eng Engine, post *redditproto.Link)
}

// InboxHandler defines methods for bots that handle new messages to their
// inbox. These include post replies, comment replies, private messages, and
// username mentions.
type InboxHandler interface {
	// Message is called when the bot receives a new private message to its
	// account. [Called as goroutine.]
	Message(eng Engine, msg *redditproto.Message)
	// Reply is called when the bot receives a reply to one of its
	// submissions in its inbox. [Called as goroutine.]
	PostReply(eng Engine, reply *redditproto.Message)
	// CommentReply is called when the bot receives a reply to one of its
	// comments in it its inbox. [Called as goroutine.]
	CommentReply(eng Engine, reply *redditproto.Message)
	// Mention is called when the bot receives a username mention in its
	// inbox. These will only appear in the user's inbox if mention
	// monitoring is turned on in preferences. [Called as goroutine.]
	Mention(eng Engine, mention *redditproto.Message)
}
