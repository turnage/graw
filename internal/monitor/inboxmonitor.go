package monitor

import (
	"github.com/turnage/graw/internal/monitor/internal/handlers"
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

const (
	// mentionSubject is the value of the subject field for comments that
	// arrive in the inbox because they mention the bot's username.
	mentionSubject = "username mention"
	// postReplySubject is the value of the subject field for comments that
	// arrive in the inbox because they are a reply to one of the bot's
	// posts.
	postReplySubject = "post reply"
)

// inboxMonitor queries a Reddit user inbox, and feeds new inbox items to the
// bot's handlers.
type inboxMonitor struct {
	base
	// postReplyHandler is the bot's interface for handling post replies.
	postReplyHandler handlers.PostReplyHandler
	// commentReplyHandler is the bot's interface for handling comment
	// replies.
	commentReplyHandler handlers.CommentReplyHandler
	// mentionHandler is the bot's interface for handling username mentions.
	mentionHandler handlers.MentionHandler
}

// InboxMonitor returns an inbox monitor that forwards events that a bot can
// handle to bot.If the bot cannot handle any inbox events, returns nil.
func InboxMonitor(
	op operator.Operator,
	messageHandler handlers.MessageHandler,
	postReplyHandler handlers.PostReplyHandler,
	commentReplyHandler handlers.CommentReplyHandler,
	mentionHandler handlers.MentionHandler,
	dir Direction,
) (Monitor, error) {
	i := &inboxMonitor{
		base: base{
			handleMessage: messageHandler.Message,
			dir:           dir,
			tip:           []string{""},
			path:          "/message/inbox",
		},
		postReplyHandler:    postReplyHandler,
		commentReplyHandler: commentReplyHandler,
		mentionHandler:      mentionHandler,
	}

	if dir == Forward {
		if err := i.sync(op); err != nil {
			return nil, err
		}
	}

	i.handleComment = i.commentDispatch
	return i, nil
}

// commentDispatch dispatches comment to their appropriate handlers, since
// comments of multiple contexts end up in the inbox.This is called as a
// goroutine from the base monitor dispatch method.
func (i *inboxMonitor) commentDispatch(comment *redditproto.Comment) {
	if comment.GetSubject() == "username mention" {
		if i.mentionHandler != nil {
			i.mentionHandler.Mention(comment)
		}
	} else if comment.GetSubject() == "post reply" {
		if i.postReplyHandler != nil {
			i.postReplyHandler.PostReply(comment)
		}
	} else {
		if i.commentReplyHandler != nil {
			i.commentReplyHandler.CommentReply(comment)
		}
	}
}
