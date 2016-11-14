// This file is a painful surrender to annoying mechanics of interfaces. These
// structs are proxies so that the interfaces satisfied by users can also be
// used by internal packages (which is the whole point).
//
// For more information, see internal/handlers.
package graw

import (
	"github.com/turnage/graw/internal/data"
)

type subredditHandlerProxy struct {
	sh SubredditHandler
}

func (s *subredditHandlerProxy) Post(p *data.Post) error {
	return s.sh.Post((*Post)(p))
}

type userHandlerProxy struct {
	uh UserHandler
}

func (u *userHandlerProxy) UserPost(p *data.Post) error {
	return u.uh.UserPost((*Post)(p))
}

func (u *userHandlerProxy) UserComment(c *data.Comment) error {
	return u.uh.UserComment((*Comment)(c))
}

type inboxHandlerProxy struct {
	ih InboxHandler
}

func (i *inboxHandlerProxy) Message(m *data.Message) error {
	return i.ih.Message((*Message)(m))
}

func (i *inboxHandlerProxy) Mention(m *data.Message) error {
	return i.ih.Mention((*Message)(m))
}

func (i *inboxHandlerProxy) CommentReply(m *data.Message) error {
	return i.ih.CommentReply((*Message)(m))
}

func (i *inboxHandlerProxy) PostReply(m *data.Message) error {
	return i.ih.PostReply((*Message)(m))
}
