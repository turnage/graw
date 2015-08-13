package graw

import (
	"github.com/turnage/redditproto"
)

// Engine defines the behaviors of the bot engine, which bots will use to
// request actions from it.
type Engine interface {}

// Run continuously runs bot, generating events from subreddits, interacting
// with reddit using agent. An error is returned if startup fails, otherwise
// this function will run until the program dies.
func Run(agent redditproto.UserAgent, bot Bot, subreddits ...string) error {
	return nil
}
