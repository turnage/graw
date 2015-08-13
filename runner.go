package graw

import (
	"fmt"
	"os"

	"github.com/turnage/redditproto"
)

const (
	// oauthURL is the url of reddit's oauth authorization server.
	oauthURL = "https://www.reddit.com/api/v1/access_token"
)

// Run continuously runs bot, generating events from subreddits, interacting
// with reddit using agent. An error is returned if startup fails, otherwise
// this function will run until the program dies.
func Run(agent *redditproto.UserAgent, bot Bot, subreddits ...string) error {
	if agent == nil {
		return fmt.Errorf("user agent was nil")
	}

	if bot == nil {
		return fmt.Errorf("bot implementation was nil")
	}

	if len(subreddits) == 0 {
		return fmt.Errorf("have no subreddits to run bot against")
	}

	httpCli, err := oauth(
		agent.GetClientId(),
		agent.GetClientSecret(),
		agent.GetUsername(),
		agent.GetPassword(),
		oauthURL)
	if err != nil {
		return err
	}

	monitor := &subredditMonitor{
		Cli:        &netClient{Client: httpCli},
		Posts:      make(chan *redditproto.Link),
		Errors:     make(chan error),
		Subreddits: subreddits,
		Kill:       make(chan bool),
	}
	go monitor.Run()

	for true {
		select {
		case post := <-monitor.Posts:
			bot.NewPost(nil, post)
		case err := <-monitor.Errors:
			fmt.Printf("error: %v.\n", err)
			os.Exit(-1)
		}
	}

	return nil
}
