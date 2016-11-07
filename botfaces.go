package graw

// SubredditHandler handles events from a subreddit.
// All methods of this interface must be goroutine safe. An error returned from
// any method of this interface will stop the graw engine.
type SubredditHandler interface {
	// Post handles a new post made in a watched subreddit.
	Post(p *Post) error
}

// UserHandler handles events from a user.
// All methods of this interface must be goroutine safe. An error returned from
// any method of this interface will stop the graw engine.
type UserHandler interface {
	// UserPost handles a new post made by a watched user.
	UserPost(p *Post) error
	// UserComment handles a new comment made by a watched user.
	UserComment(c *Comment) error
}

// InboxHandler handles events from a Reddit account inbox.
// All methods of this interface must be goroutine safe. An error returned from
// any method of this interface will stop the graw engine.
type InboxHandler interface {
	// Message handles a received private message to the bot's inbox.
	Message(m *Message) error
	// PostReply handles a reply to a bot's post received in its inbox.
	PostReply(m *Message) error
	// CommentReply handles a reply to a bot's comment received in its
	// inbox.
	CommentReply(m *Message) error
	// Mention handles a mention of the bot's username received in its
	// inbox.
	Mention(m *Message) error
}
