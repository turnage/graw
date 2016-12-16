package graw

import (
	"fmt"

	"github.com/turnage/graw/botfaces"
	"github.com/turnage/graw/reddit"
	"github.com/turnage/graw/streams"
)

var (
	postHandlerErr = fmt.Errorf(
		"You must implement PostHandler to handle subreddit feeds.",
	)
	commentHandlerErr = fmt.Errorf(
		"You must implement CommentHandler to handle subreddit " +
			"comment feeds.",
	)
	userHandlerErr = fmt.Errorf(
		"You must implement UserHandler to handle user feeds.",
	)
	loggedOutErr = fmt.Errorf(
		"You must be running as a logged in bot to get inbox feeds.",
	)
)

// Scan connects any requested logged-out event sources to the given handler,
// making requests with the given script handle. It launches a goroutine for the
// scan. It returns two functions: a stop() function to stop the scan at any
// time, and a wait() function to block until the scan fails.
func Scan(handler interface{}, script reddit.Script, cfg Config) (
	func(),
	func() error,
	error,
) {
	kill := make(chan bool)
	errs := make(chan error)

	if cfg.PostReplies || cfg.CommentReplies || cfg.Mentions || cfg.Messages {
		return nil, nil, loggedOutErr
	}

	if err := connectScanStreams(
		handler,
		script,
		cfg,
		kill,
		errs,
	); err != nil {
		return nil, nil, err
	}

	return launch(handler, kill, errs, logger(cfg.Logger))
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
		ph, ok := handler.(botfaces.PostHandler)
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

	if len(c.SubredditComments) > 0 {
		ch, ok := handler.(botfaces.CommentHandler)
		if !ok {
			return commentHandlerErr
		}

		if comments, err := streams.SubredditComments(
			sc,
			kill,
			errs,
			c.SubredditComments...,
		); err != nil {
			return err
		} else {
			go func() {
				for c := range comments {
					errs <- ch.Comment(c)
				}
			}()
		}
	}

	if len(c.Users) > 0 {
		uh, ok := handler.(botfaces.UserHandler)
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
