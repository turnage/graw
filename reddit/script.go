package reddit

import (
	"time"
)

// Script defines the behaviors of a logged out Reddit script.
type Script interface {
	Lurker
	Scanner
}

type script struct {
	Lurker
	Scanner
}

// NewScript returns a Script handle to Reddit's API which always sends the
// given agent in the user-agent header of its requests and makes requests with
// no less time between them than rate. The minimum respected value of rate is 2
// seconds, because Reddit's API rules cap logged out non-OAuth clients at 30
// requests per minute.
func NewScript(agent string, rate time.Duration) (Script, error) {
	c, err := newClient(clientConfig{agent: agent})
	r := newReaper(
		reaperConfig{
			client:     c,
			parser:     newParser(),
			hostname:   "reddit.com",
			reapSuffix: ".json",
			tls:        true,
			rate:       maxOf(rate, 2*time.Second),
		},
	)
	return &script{
		Lurker:  newLurker(r),
		Scanner: newScanner(r),
	}, err
}
