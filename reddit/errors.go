package reddit

import (
	"fmt"
)

var (
	PermissionDeniedErr   = fmt.Errorf("unauthorized access to endpoint")
	BusyErr               = fmt.Errorf("Reddit is busy right now")
	RateLimitErr          = fmt.Errorf("Reddit is rate limiting requests")
	GatewayErr            = fmt.Errorf("502 bad gateway code from Reddit")
	GatewayTimeoutErr     = fmt.Errorf("504 gateway timeout from Reddit")
	ThreadDoesNotExistErr = fmt.Errorf("The requested post does not exist.")
)
