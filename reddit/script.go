package reddit

import (
	"net/http"
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

type ScriptConfig struct {
	// Agent is the user-agent sent in all requests the bot makes through
	// this package.
	Agent string
	// Rate is the minimum amount of time between requests.
	Rate time.Duration
	// Custom HTTP client
	Client *http.Client
}

// NewScript returns a Script handle to Reddit's API which always sends the
// given agent in the user-agent header of its requests and makes requests with
// no less time between them than rate. The minimum respected value of rate is 2
// seconds, because Reddit's API rules cap logged out non-OAuth clients at 30
// requests per minute.
func NewScript(agent string, rate time.Duration) (Script, error) {
	return NewScriptFromConfig(ScriptConfig{
		Agent:  agent,
		Rate:   rate,
		Client: nil, // uses default if nil
	})
}

// NewScriptFromConfig returns a Script handle to Reddit's API from ScriptConfig
func NewScriptFromConfig(config ScriptConfig) (Script, error) {
	c, err := newClient(clientConfig{agent: config.Agent, client: config.Client})
	r := newReaper(
		reaperConfig{
			client:     c,
			parser:     newParser(),
			hostname:   "reddit.com",
			reapSuffix: ".json",
			tls:        true,
			rate:       maxOf(config.Rate, 2*time.Second),
		},
	)
	return &script{
		Lurker:  newLurker(r),
		Scanner: newScanner(r),
	}, err
}
