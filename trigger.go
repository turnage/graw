package graw

// Trigger defines the behavior of a something that can change in state without
// the runner's involvement. E.g. A trigger set up for new posts in a subreddit
// will be pulled when a new post is made.
type Trigger interface {
	// Pulled returns whether the trigger has been pulled since the last
	// call, or since instantiation if it has never been called.
	Pulled(cli client) (bool, error)
}
