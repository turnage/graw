package graw

import (
	"github.com/turnage/redditproto"
)

// Engine exposes functions of the Engine to the bot.
type Engine interface {
	// Reply posts a reply to something on reddit. The behavior depends on
	// what is being replied to. For
	//
	//   messages, this sends a private message reply.
	//   posts, this posts a top level comment.
	//   comments, this posts a comment reply.
	//
	// Use GetName() on the parent post, message, or comment to find its
	// name.
	Reply(parentName, text string) error

	// SendMessage sends a private message to a user.
	SendMessage(user, subject, text string) error

	// SelfPost makes a text post to a subreddit.
	SelfPost(subreddit, title, text string) error

	// LinkPost makes a link post to a subreddit.
	LinkPost(subreddit, title, url string) error

	// DigestThread returns a post with a parsed comment tree. Call
	// GetComments() on the returned Link for a slice of top level comments,
	// and GetReplyTree() on Comments for their replies.
	//
	// Use GetName() on a Link to find its name.
	DigestThread(threadName string) (*redditproto.Link, error)

	// Stop stops the engine. If it implemented it, the bot's TearDown
	// method will be called.
	Stop()
}

// Actor defines methods for bots that do things (send messages, make posts,
// fetch threads, etc).
type Actor interface {
	// TakeEngine is called when the engine starts; bots should save the
	// engine so they can call its methods. This is only called once.
	TakeEngine(eng Engine)
}

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
