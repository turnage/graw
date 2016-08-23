// Package grerr (graw-error) defines error values graw can encounter so that
// bots can understand and define scenarios for known error types.
package grerr

import (
	"fmt"
)

var (
	// PermissionDenied usually occurs because a bot provided invalid
	// credentials or tried to access a private subreddit.
	PermissionDenied = fmt.Errorf("reddit returned 403; permission denied")

	// Busy usually occurs when Reddit is under heavy load and does not have
	// anything to do with the running bot.
	Busy = fmt.Errorf("reddit returned 503; it is busy right now")

	// RateLimit means Reddit has received too many requests in the allowed
	// interval for the bot's user agent. This usually means the bot did not
	// correctly define or is using a default user agent, or is running on
	// multiple instances of graw, because graw automatically enforces rule
	// abiding rate limits on all bots.
	RateLimit = fmt.Errorf("reddit returned 429; too many requests")
)
