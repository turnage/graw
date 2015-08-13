package graw

import (
	"github.com/turnage/redditproto"
)

// Run continuously runs bot, generating events from subreddits, interacting
// with reddit using agent. An error is returned if startup fails, otherwise
// this function will run until the program dies.
func Run(agent redditproto.UserAgent, bot Bot, subreddits ...string) error {
	return nil
}
