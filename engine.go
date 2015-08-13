package graw

import (
	"fmt"

	"github.com/turnage/redditproto"
)

// These refresh rates define which event generators (e.g. subreddit monitor)
// get how many queries per minute. A single bot, according to Reddit's api
// rules, can make 60 queries per minute. These must all sum to 60 or less.
const (
	newPostRefreshRate = 30
)

// Engine defines the behaviors of the bot engine, which bots will use to
// request actions from it.
type Engine interface{}

// engine implements Engine.
type engine struct {
	// bot is the bot this engine runs.
	bot Bot
	// cli is the clinet engine uses to communicated with reddit.
	cli client
	// errors is the channel the engine's goroutines will report their
	// errors to.
	errors chan error
	// monitor watches subreddits for new posts.
	monitor subredditMonitor
}

// makeEngine returns an engine that uses agent to authenticate with reddit, and
// runs bot against subreddits.
func makeEngine(agent *redditproto.UserAgent, bot Bot, subreddits []string) (*engine, error) {
	if agent == nil {
		return nil, fmt.Errorf("user agent was nil")
	}

	if bot == nil {
		return nil, fmt.Errorf("bot implementation was nil")
	}

	if len(subreddits) == 0 {
		return nil, fmt.Errorf("have no subreddits to run bot against")
	}

	httpCli, err := oauth(
		agent.GetClientId(),
		agent.GetClientSecret(),
		agent.GetUsername(),
		agent.GetPassword(),
		oauthURL)
	if err != nil {
		return nil, err
	}

	errors := make(chan error)
	return &engine{
		bot:    bot,
		cli:    &netClient{Client: httpCli},
		errors: errors,
		monitor: subredditMonitor{
			Posts:       make(chan *redditproto.Link),
			Errors:      errors,
			Subreddits:  subreddits,
			Kill:        make(chan bool),
			RefreshRate: newPostRefreshRate,
		},
	}, nil
}

// Run runs the engine until an error occurs or its bot stops it.
func (e *engine) Run() error {
	go e.monitor.Run(e.cli)
	for true {
		select {
		case post := <-e.monitor.Posts:
			e.bot.NewPost(e, post)
		case err := <-e.errors:
			e.monitor.Kill <- true
			return err
		}
	}
	return nil
}
