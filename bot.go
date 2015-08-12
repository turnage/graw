package graw

import (
	"github.com/turnage/redditproto"
)

// Bot defines the expecations the runner has of a bot. Each method is an event
// that a Bot implementation is expected to handle.
//
// These methods WILL be run concurrently, and should be implemented with that
// as the expected use case.
type Bot interface {
	// When the runner finds a new post in a watched subreddit, it will call
	// this method on the bot for the new post. This will be called on every
	// new post found. "New"  in this context means "posted after the runner
	// started".
	NewPost(user User, post *redditproto.Link) error
}
