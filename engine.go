package graw

import (
	"github.com/turnage/redditproto"
)

// Engine exposes certain functions to the bot the engine is running.
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
