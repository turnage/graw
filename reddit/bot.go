package reddit

import (
	"time"
)

type BotConfig struct {
	Agent string
	App   App
	Rate  time.Duration
}

type Bot interface {
	Account
	Lurker
	Scanner
}

type bot struct {
	Account
	Lurker
	Scanner
}

func NewBot(c BotConfig) (Bot, error) {
	cli, err := newClient(clientConfig{agent: c.Agent, app: c.App})
	r := newReaper(
		reaperConfig{
			client:   cli,
			parser:   newParser(),
			hostname: "oauth.reddit.com",
			tls:      true,
		},
	)
	return &bot{
		Account: newAccount(r),
		Lurker:  newLurker(r),
		Scanner: newScanner(r),
	}, err
}

func NewBotFromAgentFile(filename string, rate time.Duration) (Bot, error) {
	agent, app, err := load(filename)
	if err != nil {
		return nil, err
	}

	return NewBot(
		BotConfig{
			Agent: agent,
			App:   app,
			Rate:  rate,
		},
	)
}
