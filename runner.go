package graw

import (
	"github.com/turnage/redditproto"
)

const (
	// oauthURL is the url of reddit's oauth authorization server.
	oauthURL = "https://www.reddit.com/api/v1/access_token"
	// maxQueries is the amount of queries reddit allows a bot to make per
	// minute.
	maxQueries = 60
)

// Run continuously runs bot, generating events from subreddits, interacting
// with reddit using agent. An error is returned if startup fails, otherwise
// this function will run until the program dies.
func Run(agent *redditproto.UserAgent, bot Bot, subreddits ...string) error {
	eng, err := makeEngine(agent, bot, subreddits)
	if err != nil {
		return err
	}

	return eng.Run()
}
