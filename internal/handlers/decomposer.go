package handlers

import (
	"github.com/turnage/graw/internal/data"
)

// DecomposeSubredditHandler breaks a SubredditHandler into a generic post
// handler.
func DecomposeSubredditHandler(sh SubredditHandler) PostHandler {
	return PostHandlerFunc(sh.Post)
}

// DecomposeUserHandler breaks a UserHandler into handlers for the posts and
// comments that the user makes.
func DecomposeUserHandler(uh UserHandler) (PostHandler, CommentHandler) {
	return PostHandlerFunc(uh.UserPost), CommentHandlerFunc(uh.UserComment)
}

// DecomposeInboxHandler breaks an InboxHandler down to a handler for the
// messages it receives.
func DecomposeInboxHandler(ih InboxHandler) MessageHandler {
	return MessageHandlerFunc(
		func(m *data.Message) error {
			if m.WasComment {
				switch m.Subject {
				case "comment reply":
					return ih.CommentReply(m)
				case "post reply":
					return ih.PostReply(m)
				case "username mention":
					return ih.Mention(m)
				}
			}

			return ih.Message(m)
		},
	)
}
