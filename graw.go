package graw

import (
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/operator"
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
	Post(Engine, *redditproto.Link)
	// Message will be called to handle new private messages to the Bot's
	// inbox.
	Message(Engine, *redditproto.Message)
	// Reply will be called to handle comment replies to the bot.
	Reply(Engine, *redditproto.Message)
	// TearDown will be called at the end of execution so the bot can free
	// its resources. It will not be run in parallel.
	TearDown()
}

// Run runs a bot against live reddit. agent should be the filename of a
// configured user agent protobuffer. The bot will monitor all provide
// subreddits.
//
// See the wiki for more details.
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
