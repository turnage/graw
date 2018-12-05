package reddit

import (
	"errors"
)

// Errors that can be returned by the Reddit API.
var (
	ErrPermissionDenied = errors.New("unauthorized access to endpoint")
	ErrBusy             = errors.New("Reddit is busy right now")
	ErrRateLimit        = errors.New("Reddit is rate limiting requests")
	ErrBadGateway       = errors.New("502 bad gateway code from Reddit")
	ErrGatewayTimeout   = errors.New("504 gateway timeout code from Reddit")
	ErrThreadNotExists  = errors.New("the requested thread does not exist")
)
