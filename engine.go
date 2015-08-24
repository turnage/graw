package graw

import (
	"github.com/turnage/redditproto"
)

// Engine defines the interface for bots to interact with the engine. These
// methods are requests to the engine to perform actions on behalf of the bot,
// when it decides it is time.
type Engine interface {
	// ReplyToInbox sends a reply to something in the inbox. This posts a
	// comment when replying to a comment or post, and sends a message when
	// replying to a private message. Non-message types will be
	// represented as a Message when they appear in the inbox (such as a
	// self post that mentions the bot's username).
	ReplyToInbox(msg *redditproto.Message, text string) error
	// SendMessage sends a private message to a user.
	SendMessage(user, subject, text string) error

	// ReplyToPost posts a top level comment on a submission.
	ReplyToPost(post *redditproto.Link, text string) error
	// SelfPost makes a text post to a subreddit.
	SelfPost(subreddit, title, text string) error
	// LinkPost makes a link post to a subreddit.
	LinkPost(subreddit, title, url string) error
	// ScrapeThread returns the full thread for a post, including the entire
	// comment tree. Access the link's GetComments() method for the comment
	// tree, and comment subtrees using the comments' GetReplyTree() method.
	ScrapeThread(thread *redditproto.Link) (*redditproto.Link, error)
	// ScrapeThreadAt works like ScrapeThread, but takes a permalink to a
	// thread. A permalink is _not_ a full url; it is the part that starts
	// with /r/.
	ScrapeThreadAt(permalink string) (*redditproto.Link, error)

	// Stop stops the engine execution.
	Stop()
}
