package graw

import (
	"fmt"

	"github.com/turnage/graw/reddit"
	"github.com/turnage/graw/streams"
)

var (
	postHandlerErr = fmt.Errorf(
		"You must implement PostHandler to handle subreddit feeds.",
	)
	userHandlerErr = fmt.Errorf(
		"You must implement UserHandler to handle user feeds.",
	)
	loggedOutErr = fmt.Errorf(
		"You must be running as a logged in bot to get inbox feeds.",
	)
)

func Scan(handler interface{}, script reddit.Script, cfg Config) error {
	kill := make(chan bool)
	errs := make(chan error)

	if cfg.PostReplies || cfg.CommentReplies || cfg.Mentions || cfg.Messages {
		return loggedOutErr
	}

	if err := connectScanStreams(
		handler,
		script,
		cfg,
		kill,
		errs,
	); err != nil {
		return err
	}

	return foreman(kill, errs, logger(cfg.Logger))
}

// connectScanStreams connects the streams a scanner can subscribe to to the
// handler.
func connectScanStreams(
	handler interface{},
	sc reddit.Scanner,
	c Config,
	kill <-chan bool,
	errs chan<- error,
) error {
	if len(c.Subreddits) > 0 {
		ph, ok := handler.(PostHandler)
		if !ok {
			return postHandlerErr
		}

		if posts, err := streams.Subreddits(
			sc,
			kill,
			errs,
			c.Subreddits...,
		); err != nil {
			return err
		} else {
			go func() {
				for p := range posts {
					errs <- ph.Post(p)
				}
			}()
		}
	}

	if len(c.Users) > 0 {
		uh, ok := handler.(UserHandler)
		if !ok {
			return userHandlerErr
		}

		for _, user := range c.Users {
			if posts, comments, err := streams.User(
				sc,
				kill,
				errs,
				user,
			); err != nil {
				return err
			} else {
				go func() {
					for p := range posts {
						errs <- uh.UserPost(p)
					}
				}()
				go func() {
					for c := range comments {
						errs <- uh.UserComment(c)
					}
				}()
			}
		}
	}

	return nil
}
