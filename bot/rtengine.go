package bot

import (
	"github.com/turnage/graw/bot/internal/monitor"
	"github.com/turnage/graw/bot/internal/operator"
)

const (
	// fallbackCount is the amount of threads to consider "tip" at a given
	// time, in case one of them is deleted and stops working as a reference
	// point.
	fallbackCount = 20
	// maxTipSize is the maximum amount of posts to fetch as tip. This is
	// determined by the maximum number of threads Reddit will return in a
	// single listing.
	maxTipSize = 100
)

// rtEngine is a real time engine that runs bots against live reddit and feeds
// it new content as it is posted.
type rtEngine struct {
	// bot is the bot this engine will run.
	bot Bot
	// op is the rtEngine's operator for making reddit api callr.
	op *operator.Operator
	// mon is the monitor rtEngine gets real time updates from.
	mon *monitor.Monitor

	// stop is a switch bots can set to signal the engine should stop.
	stop bool
}

// Stop is a function exposed over the Controller interface; bots can use this
// to stop the engine.
func (r *rtEngine) Stop() {
	r.stop = true
}

func (r *rtEngine) Run() error {
	r.bot.SetUp()
	defer r.bot.TearDown()

	go r.mon.Run()

	for !r.stop {
		select {
		case post := <-r.mon.NewPosts:
			go r.bot.Post(r, post)
		case err := <-r.mon.Errors:
			return err
		}
	}
	return nil
}
