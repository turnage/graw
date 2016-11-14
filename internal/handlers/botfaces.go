// This file is a mirror of graw/botfaces.go. It is unfortunate that this has to
// exist. Here is why:
//
// I don't need to export these interfaces to users for users to implement them,
// however to be implemented they must receive the data types. I export the data
// types at the top level package with a typedef, but interfaces cannot be
// satisfied by receiving a typedef of the type in the interface.
//
// So, if I want to have access to the interfaces inside of graw, I need to
// duplicate them here, and create a converter function at the top level which
// has access to the exported types.
package handlers

import (
	"github.com/turnage/graw/internal/data"
)

type SubredditHandler interface {
	Post(p *data.Post) error
}

type UserHandler interface {
	UserPost(p *data.Post) error
	UserComment(c *data.Comment) error
}

type InboxHandler interface {
	Message(m *data.Message) error
	PostReply(m *data.Message) error
	CommentReply(m *data.Message) error
	Mention(m *data.Message) error
}
