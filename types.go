package graw

import (
	"github.com/turnage/graw/internal/data"
)

// Post represents self and link posts on Reddit.
type Post data.Post

// Comment represents a comment on Reddit.
type Comment data.Comment

// Message represents a Reddit element in the bot's inbox (comment replies to
// the bot's posts or comments will come as Messages in the inbox).
type Message data.Message
