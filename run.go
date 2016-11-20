package graw

import (
	"fmt"

	"github.com/turnage/graw/botfaces"
	"github.com/turnage/graw/reddit"
	"github.com/turnage/graw/streams"
)

var (
	postReplyHandlerErr = fmt.Errorf(
		"You must implement PostReplHandler to take post reply feeds.",
	)
	commentReplyHandlerErr = fmt.Errorf(
		"You must implement CommentReplyHandler to take comment reply feeds.",
	)
	mentionHandlerErr = fmt.Errorf(
		"You must implement MentionHandler to take mention feeds.",
	)
	messageHandlerErr = fmt.Errorf(
		"You must implement MessageHandler to take message feeds.",
	)
)

// Run connects a handler to any requested event sources and makes requests with
// the given bot api handle. It launches a goroutine for the run. It returns two
// functions, a stop() function to terminate the graw run at any time, and a
// wait() function to block until the graw run fails.
func Run(handler interface{}, bot reddit.Bot, cfg Config) (
	func(),
	func() error,
	error,
) {
	kill := make(chan bool)
	errs := make(chan error)

	if err := connectAllStreams(
		handler,
		bot,
		cfg,
		kill,
		errs,
	); err != nil {
		return nil, nil, err
	}

	return launch(handler, kill, errs, logger(cfg.Logger))
}

func connectAllStreams(
	handler interface{},
	bot reddit.Bot,
	c Config,
	kill <-chan bool,
	errs chan<- error,
) error {
	if err := connectScanStreams(
		handler,
		bot,
		c,
		kill,
		errs,
	); err != nil {
		return err
	}

	// lol no generics:

	if c.PostReplies {
		if prh, ok := handler.(botfaces.PostReplyHandler); !ok {
			return postReplyHandlerErr
		} else if prs, err := streams.PostReplies(
			bot,
			kill,
			errs,
		); err != nil {
			return err
		} else {
			go func() {
				for pr := range prs {
					errs <- prh.PostReply(pr)
				}
			}()
		}
	}

	if c.CommentReplies {
		if crh, ok := handler.(botfaces.CommentReplyHandler); !ok {
			return commentReplyHandlerErr
		} else if crs, err := streams.CommentReplies(
			bot,
			kill,
			errs,
		); err != nil {
			return err
		} else {
			go func() {
				for cr := range crs {
					errs <- crh.CommentReply(cr)
				}
			}()
		}
	}

	if c.Mentions {
		if mh, ok := handler.(botfaces.MentionHandler); !ok {
			return mentionHandlerErr
		} else if ms, err := streams.Mentions(
			bot,
			kill,
			errs,
		); err != nil {
			return err
		} else {
			go func() {
				for m := range ms {
					errs <- mh.Mention(m)
				}
			}()
		}
	}

	if c.Messages {
		if mh, ok := handler.(botfaces.MessageHandler); !ok {
			return messageHandlerErr
		} else if ms, err := streams.Messages(
			bot,
			kill,
			errs,
		); err != nil {
			return err
		} else {
			go func() {
				for m := range ms {
					errs <- mh.Message(m)
				}
			}()
		}
	}

	return nil
}
