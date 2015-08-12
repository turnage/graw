package graw

import (
	"github.com/turnage/redditproto"
)

// Run continuously runs bot, generating events from triggers, interacting with
// reddit using agent. An error is returned if startup fails, otherwise this
// function will run until the program dies.
func Run(agent redditproto.UserAgent, bot Bot, triggers []Trigger) error {
	return nil
}
