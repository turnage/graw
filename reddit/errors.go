package reddit

import (
	"errors"
)

// Errors that can be returned by the Reddit API.
var (
	PermissionDeniedErr   = errors.New("unauthorized access to endpoint")
	BusyErr               = errors.New("Reddit is busy right now")
	RateLimitErr          = errors.New("Reddit is rate limiting requests")
	GatewayErr            = errors.New("502 bad gateway code from Reddit")
	GatewayTimeoutErr     = errors.New("504 gateway timeout code from Reddit")
	ThreadDoesNotExistErr = errors.New("the requested post does not exist")
)
