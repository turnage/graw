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
	// Bot is the InboxHandler InboxMonitor will send the inbox items it
	// finds to.
	Bot api.InboxHandler
}

// Update updates the inbox and sends all messages to the InboxHandler.
func (i *InboxMonitor) Update() error {
	messages, err := i.Op.Inbox()
	if err != nil {
		return err
	}

	for _, message := range messages {
		if message.GetSubject() == "username mention" {
			go i.Bot.Mention(message)
		} else if message.GetSubject() == "post reply" {
			go i.Bot.PostReply(message)
		} else if message.GetWasComment() {
			go i.Bot.CommentReply(message)
		} else {
			go i.Bot.Message(message)
		}
	}

	return nil
}
