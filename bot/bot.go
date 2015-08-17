package bot

import (
	"github.com/turnage/graw/bot/internal/client"
	"github.com/turnage/graw/bot/internal/operator"
	"github.com/turnage/redditproto"
)

// Bot defines the behaviors of a bot graw will run.
//
// The graw engine will generate events, and call Bot's methods to handle them.
// Bot implementations should expect that their methods will be called as
// goroutines, and be safe to run concurrently.
type Bot interface {
	// Post will be called to handle events that yield a post the Bot has
	// not seen before.
	Post(contr Controller, post *redditproto.Link) error
}

// Run runs a bot against live reddit. agent should be the filename of an
// authenticated user agent (see "graw grant"). Events will be generated from
// all included subreddits.
func Run(agent string, bot Bot, subreddits ...string) error {
	cli, err := client.New(agent)
	if err != nil {
		return err
	}

	eng := &rtEngine{
		bot:        bot,
		op:         operator.New(cli),
		subreddits: subreddits,
	}
	return eng.Run()
}
