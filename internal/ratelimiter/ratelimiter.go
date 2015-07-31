// Package ratelimiter provides rate limiters which control the rate at which
// actions are performed.
package ratelimiter

// RateLimiter defines the behavior of something that limits the rate of
// requests. They govern the rate of an action by permitting or denying access
// to that action in a given instance.
type RateLimiter interface {
	// Request requests access to perform the rate limited action. If block
	// is true, Request just waits for the action to be allowed and always
	// returns true. If block is false, returns true or false to indicate
	// whether the action is allowed.
	Request(block bool) bool
}
