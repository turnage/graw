package monitor

import (
	"github.com/turnage/graw/api"
	"github.com/turnage/graw/internal/operator"
)

// InboxMonitor queries a Reddit user inbox, and feeds new inbox items to its
// InboxHandler.
type InboxMonitor struct {
	// Op is the operator InboxMonitor will use to query Reddit.
	Op operator.Operator
	// MessageHandler is the bot's interface for handling messages.
	MessageHandler api.MessageHandler
	// PostReplyHandler is the bot's interface for handling post replies.
	PostReplyHandler api.PostReplyHandler
	// CommentReplyHandler is the bot's interface for handling comment
	// replies.
	CommentReplyHandler api.CommentReplyHandler
	// MentionHandler is the bot's interface for handling username mentions.
	MentionHandler api.MentionHandler
}

// Update updates the inbox and sends all inbox items to their handler.
func (i *InboxMonitor) Update() error {
	messages, err := i.Op.Inbox()
	if err != nil {
		return err
	}

	handled := []string{}
	for _, message := range messages {
		if message.GetSubject() == "username mention" {
			if i.MentionHandler != nil {
				go i.MentionHandler.Mention(message)
				handled = append(handled, message.GetName())
			}
		} else if message.GetSubject() == "post reply" {
			if i.PostReplyHandler != nil {
				go i.PostReplyHandler.PostReply(message)
				handled = append(handled, message.GetName())
			}
		} else if message.GetWasComment() {
			if i.CommentReplyHandler != nil {
				go i.CommentReplyHandler.CommentReply(message)
				handled = append(handled, message.GetName())
			}
		} else {
			if i.MessageHandler != nil {
				go i.MessageHandler.Message(message)
				handled = append(handled, message.GetName())
			}
		}
	}

	return i.Op.MarkAsRead(handled...)
}
