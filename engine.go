package graw

import (
	"github.com/turnage/redditproto"
)

// Engine defines the interface for bots to interact with the engine. These
// methods are requests to the engine to perform actions on behalf of the bot,
// when it decides it is time.
type Engine interface {
	// ReplyToPost posts a top level comment on a submission.
	ReplyToPost(post *redditproto.Link, text string) error
	// ReplyToInbox sends a reply to something in the inbox. This posts a
	// comment when replying to a comment reply, and sends a message when
	// replying to a private message.
	ReplyToInbox(msg *redditproto.Message, text string) error
	// Stop stops the engine execution.
	Stop()
}
