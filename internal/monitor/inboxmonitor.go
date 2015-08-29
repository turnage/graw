package monitor

import (
	"github.com/turnage/graw/api"
	"github.com/turnage/graw/internal/operator"
)

// inboxMonitor queries a Reddit user inbox, and feeds new inbox items to the
// bot's handlers.
type inboxMonitor struct {
	// op is the operator inboxMonitor will use to query Reddit.
	op operator.Operator
	// messageHandler is the bot's interface for handling messages.
	messageHandler api.MessageHandler
	// postReplyHandler is the bot's interface for handling post replies.
	postReplyHandler api.PostReplyHandler
	// commentReplyHandler is the bot's interface for handling comment
	// replies.
	commentReplyHandler api.CommentReplyHandler
	// mentionHandler is the bot's interface for handling username mentions.
	mentionHandler api.MentionHandler
}

// InboxMonitor returns an inbox monitor that forwards events that a bot can
// handle to bot. If the bot cannot handle any inbox events, returns nil.
func InboxMonitor(op operator.Operator, bot interface{}) Monitor {
	messageHandler, _ := bot.(api.MessageHandler)
	postReplyHandler, _ := bot.(api.PostReplyHandler)
	commentReplyHandler, _ := bot.(api.CommentReplyHandler)
	mentionHandler, _ := bot.(api.MentionHandler)

	if messageHandler == nil &&
		postReplyHandler == nil &&
		commentReplyHandler == nil &&
		mentionHandler == nil {
		return nil
	}

	return &inboxMonitor{
		op:                  op,
		messageHandler:      messageHandler,
		postReplyHandler:    postReplyHandler,
		commentReplyHandler: commentReplyHandler,
		mentionHandler:      mentionHandler,
	}
}

// Update updates the inbox and sends all inbox items to their handler.
func (i *inboxMonitor) Update() error {
	messages, err := i.op.Inbox()
	if err != nil {
		return err
	}

	handled := []string{}
	for _, message := range messages {
		if message.GetSubject() == "username mention" {
			if i.mentionHandler != nil {
				go i.mentionHandler.Mention(message)
				handled = append(handled, message.GetName())
			}
		} else if message.GetSubject() == "post reply" {
			if i.postReplyHandler != nil {
				go i.postReplyHandler.PostReply(message)
				handled = append(handled, message.GetName())
			}
		} else if message.GetWasComment() {
			if i.commentReplyHandler != nil {
				go i.commentReplyHandler.CommentReply(message)
				handled = append(handled, message.GetName())
			}
		} else {
			if i.messageHandler != nil {
				go i.messageHandler.Message(message)
				handled = append(handled, message.GetName())
			}
		}
	}

	if len(handled) > 0 {
		return i.op.MarkAsRead(handled...)
	}

	return nil
}
