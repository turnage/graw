package ratelimiter

import (
	"testing"
	"time"
)

func TestNewTimeRateLimiter(t *testing.T) {
	limiter := NewTimeRateLimiter(time.Second, 1).(*timeRateLimiter)
	if limiter.windowSize != time.Second {
		t.Error("window size not set")
	}
	if limiter.actionsPerWindow != 1 {
		t.Error("actions per window not set")
	}
}

func TestTimeRateLimiterRequest(t *testing.T) {
	limiter := timeRateLimiter{windowSize: time.Second, actionsPerWindow: 1}
	if !limiter.Request(false) {
		t.Error("should have allowed first action in window")
	}
	if limiter.Request(false) {
		t.Error("should have denied second action in window")
	}
	if !limiter.Request(true) {
		t.Error("should have blocked until third action was allowed")
	}
}
