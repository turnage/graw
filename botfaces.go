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
