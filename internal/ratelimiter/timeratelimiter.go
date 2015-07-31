package ratelimiter

import (
	"time"
)

type timeRateLimiter struct {
	actionsPerWindow  uint
	actionsThisWindow uint
	windowStart       time.Time
	windowSize        time.Duration
}

// NewTimeRateLimiter returns a new RateLimiter implementation that allows a
// maximum number of actions within a window of time.
//
// NetTimeRateLimiter(time.Minute, 60) will return a RateLimiter that only
// allows an action 60 times per minute, and cuts off access until the next
// minute if 60 have been performed.
func NewTimeRateLimiter(window time.Duration, actions uint) RateLimiter {
	return &timeRateLimiter{windowSize: window, actionsPerWindow: actions}
}

func (t *timeRateLimiter) Request(block bool) bool {
	if t.windowStart.IsZero() || time.Since(t.windowStart) > t.windowSize {
		t.windowStart = time.Now()
		t.actionsThisWindow = 0
	}

	if t.actionsThisWindow < t.actionsPerWindow {
		t.actionsThisWindow++
		return true
	}

	if block {
		time.Sleep(t.windowStart.Add(t.windowSize).Sub(time.Now()))
	} else {
		return false
	}

	return true
}
