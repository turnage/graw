package reddit

import (
	"time"
)

type Script interface {
	Lurker
	Scanner
}

type script struct {
	Lurker
	Scanner
}

func NewScript(agent string, rate time.Duration) (Script, error) {
	c, err := newClient(clientConfig{agent: agent})
	r := newReaper(
		reaperConfig{
			client:   c,
			parser:   newParser(),
			hostname: "reddit.com",
			tls:      true,
		},
	)
	return &script{
		Lurker:  newLurker(r),
		Scanner: newScanner(r),
	}, err
}
