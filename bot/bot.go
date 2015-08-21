package bot

import (
	"github.com/turnage/graw/bot/internal/monitor"
	"github.com/turnage/graw/bot/internal/operator"
	"github.com/turnage/redditproto"
)

// Bot defines the behaviors of a bot graw will run.
//
// The graw engine will generate events, and call Bot's methods to handle them.
//
// Bot implementations should expect that their methods will be called as
// goroutines, and be safe to run in parallel. SetUp() and TearDown() are exempt
// from this requirment.
type Bot interface {
	// SetUp will be called immediately before the start of the engine. It
	// will not be run in parallel.
	SetUp()
	// Post will be called to handle events that yield a post the Bot has
	// not seen before.
	Post(Controller, *redditproto.Link)
	// Message will be called to handle new private messages to the Bot's
	// inbox.
	Message(Controller, *redditproto.Message)
	// Reply will be called to handle comment replies to the bot.
	Reply(Controller, *redditproto.Message)
	// TearDown will be called at the end of execution so the bot can free
	// its resources. It will not be run in parallel.
	TearDown()
}

// Run runs a bot against live reddit. agent should be the filename of an
// authenticated user agent (see "graw grant"). Events will be generated from
// all included subreddits.
func Run(agent string, bot Bot, subreddits ...string) error {
	op, err := operator.New(agent)
	if err != nil {
		return err
	}

	eng := &rtEngine{
		bot: bot,
		op:  op,
		mon: monitor.New(op, subreddits),
	}

	return eng.Run()
}
