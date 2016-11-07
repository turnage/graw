package graw

import (
	"fmt"
	"time"

	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/client"
	"github.com/turnage/graw/internal/data"
	"github.com/turnage/graw/internal/dispatcher"
	"github.com/turnage/graw/internal/engine"
	"github.com/turnage/graw/internal/handlers"
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/reap"
	"github.com/turnage/graw/internal/rsort"
)

var (
	subredditHandlerErr = fmt.Errorf(
		"Config requests subreddit events, but bot does not " +
			"implement SubredditHandler interface.",
	)

	userHandlerErr = fmt.Errorf(
		"Config requests user events, but bot does not implement " +
			"UserHandler interface.",
	)
)

// minimumInterval is the minimum interval between requests a bot is allowed in
// order to comply with Reddit's API rules.
const minimumInterval = time.Second

// tokenURL is the url to request OAuth2 tokens from reddit.
const tokenURL = "https://www.reddit.com/api/v1/access_token"

// Login state parameters; these are maps of "Logged in user?" to parameter.
var (
	hostname = map[bool]string{
		true:  "oauth.reddit.com",
		false: "www.reddit.com",
	}
)

func Run(c Config, bot interface{}) error {
	reaper, isUser, err := buildReaper(c)
	if err != nil {
		return err
	}

	dispatchers := []dispatcher.Dispatcher{}

	if len(c.Subreddits) > 0 {
		if _, ok := bot.(SubredditHandler); !ok {
			return subredditHandlerErr
		}

		path := subredditsPath(c.Subreddits)
		if !isUser {
			path = logPathsOut([]string{path})[0]
		}

		if mon, err := monitor.New(
			monitor.Config{
				Path:   path,
				Lurker: api.NewLurker(reaper),
				Sorter: rsort.New(),
			},
		); err != nil {
			return err
		} else {
			handler := handlers.PostHandlerFunc(
				func(p *data.Post) error {
					return bot.(SubredditHandler).Post(
						(*Post)(p),
					)
				},
			)
			dispatchers = append(
				dispatchers, dispatcher.New(
					dispatcher.Config{
						Monitor:     mon,
						PostHandler: handler,
					},
				),
			)
		}
	}

	return engine.New(
		engine.Config{
			Dispatchers: dispatchers,
			Rate:        rateLimit(c.Rate),
		},
	).Run()
}

// rateLimit returns a rate limiter compliant with the Reddit API.
func rateLimit(interval time.Duration) <-chan time.Time {
	if interval < minimumInterval {
		interval = minimumInterval
	}
	return time.Tick(interval)
}

// buildReaper returns a reaper built with the config and whether the reaper
// acts as a logged in user.
func buildReaper(c Config) (reap.Reaper, bool, error) {
	isUser := false

	app := client.App{}
	if c.App != nil {
		isUser = true
		app = client.App{
			TokenURL: tokenURL,
			ID:       c.App.ID,
			Secret:   c.App.Secret,
			Username: c.App.Username,
			Password: c.App.Password,
		}
	}

	cli, err := client.New(client.Config{c.Agent, app})
	return reap.New(
		reap.Config{
			Client:   cli,
			Parser:   data.NewParser(),
			Hostname: hostname[isUser],
			TLS:      true,
		},
	), isUser, err
}
