// This file is a painful surrender to annoying mechanics of interfaces. These
// structs are proxies so that the interfaces satisfied by users can also be
// used by internal packages (which is the whole point).
//
// For more information, see internal/handlers.
package graw

import (
	"github.com/turnage/graw/internal/data"
	"github.com/turnage/graw/internal/handlers"
)

type subredditHandlerProxy struct {
	sh SubredditHandler
}

func (s *subredditHandlerProxy) Post(p *data.Post) error {
	return s.sh.Post((*Post)(p))
}

func subredditHandlerProxyFrom(s SubredditHandler) handlers.SubredditHandler {
	if s == nil {
		return nil
	}

	return &subredditHandlerProxy{s}
}

type userHandlerProxy struct {
	uh UserHandler
}

func userHandlerProxyFrom(u UserHandler) handlers.UserHandler {
	if u == nil {
		return nil
	}

	return &userHandlerProxy{u}
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

func inboxHandlerProxyFrom(i InboxHandler) handlers.InboxHandler {
	if i == nil {
		return nil
	}

	return &inboxHandlerProxy{i}
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
