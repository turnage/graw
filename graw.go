package graw

import (
	"log"
	"os"
	"time"

	"github.com/turnage/graw/internal/client"
	"github.com/turnage/graw/internal/data"
	"github.com/turnage/graw/internal/dispatcher"
	"github.com/turnage/graw/internal/engine"
	"github.com/turnage/graw/internal/reap"
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
	reaper, loggedIn, err := buildReaper(c)
	if err != nil {
		return err
	}

	dispatchers := []dispatcher.Dispatcher{}

	if len(c.Subreddits) > 0 {
		if d, err := subredditStream(
			c.Subreddits,
			loggedIn,
			reaper,
			bot,
		); err != nil {
			return err
		} else {
			dispatchers = append(dispatchers, d)
		}
	}

	if len(c.Users) > 0 {
		if d, err := userStreams(
			c.Users,
			loggedIn,
			reaper,
			bot,
		); err != nil {
			return err
		} else {
			dispatchers = append(dispatchers, d...)
		}
	}

	return engine.New(
		engine.Config{
			Dispatchers: dispatchers,
			Rate:        rateLimit(c.Rate, loggedIn),
			Logger:      log.New(os.Stderr, "", log.LstdFlags),
		},
	).Run()
}

// rateLimit returns a rate limiter compliant with the Reddit API.
func rateLimit(interval time.Duration, loggedIn bool) <-chan time.Time {
	minimum := minimumInterval
	if !loggedIn {
		minimum *= 2
	}
	if interval < minimum {
		interval = minimum
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
