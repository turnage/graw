// Package streams provide an event stream from Reddit. This package abstracts
// dispatcher.
package streams

import (
	"log"

	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/dispatcher"
	"github.com/turnage/graw/internal/handlers"
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/reap"
	"github.com/turnage/graw/internal/rsort"
)

// Config configures a stream set.
type Config struct {
	// LoggedIn indicates whether the bot is running as a logged in user.
	LoggedIn bool
	// Subreddits are a set of subreddits to generate events for new posts
	// in.
	Subreddits       []string
	SubredditHandler handlers.SubredditHandler
	// Users are a set of users to generate events for new posts or comments
	// they make.
	Users       []string
	UserHandler handlers.UserHandler
	// Inbox indicates whether the bot wants events from its inbox.
	Inbox        bool
	InboxHandler handlers.InboxHandler
	// Reaper is the reaper used to make requests.
	Reaper reap.Reaper
	// Bot is the bot who will be fed events for processing.
	Bot interface{}
}

// singleConfig configures one stream. All sources will be reduced to a stream
// config before building.
type singleConfig struct {
	// path is the url path to the event source.
	path string

	ph handlers.PostHandler
	ch handlers.CommentHandler
	mh handlers.MessageHandler
}

func New(c Config) ([]dispatcher.Dispatcher, error) {
	var cfgs []singleConfig

	if cfg, err := subreddits(c.Subreddits, c.SubredditHandler); err != nil {
		return nil, err
	} else {
		cfgs = append(cfgs, cfg)
	}

	if cs, err := users(c.Users, c.UserHandler); err != nil {
		return nil, err
	} else {
		cfgs = append(cfgs, cs...)
	}

	if c.Inbox {
		if cfg, err := inbox(c.LoggedIn, c.InboxHandler); err != nil {
			return nil, err
		} else {
			cfgs = append(cfgs, cfg)
		}
	}

	var streams []dispatcher.Dispatcher
	for _, sc := range cfgs {
		if sc.path == "" {
			continue
		}

		if !c.LoggedIn {
			sc.path = logPathOut(sc.path)
		}

		log.Printf("Adding cfg: %v\n", sc)

		if mon, err := monitor.New(
			monitor.Config{
				Path:   sc.path,
				Lurker: api.NewLurker(c.Reaper),
				Sorter: rsort.New(),
			},
		); err != nil {
			return nil, err
		} else {
			streams = append(
				streams, dispatcher.New(
					dispatcher.Config{
						Monitor:        mon,
						PostHandler:    sc.ph,
						CommentHandler: sc.ch,
						MessageHandler: sc.mh,
					},
				),
			)
		}

	}

	return streams, nil
}
