package graw

import (
	"fmt"

	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/data"
	"github.com/turnage/graw/internal/dispatcher"
	"github.com/turnage/graw/internal/handlers"
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/reap"
	"github.com/turnage/graw/internal/rsort"
)

var (
	subredditHandlerErr = fmt.Errorf(
		"Config requests subreddit events, but bot does not " +
			"implement SubredditHandler interface.",
	)

	userHandlerErr = fmt.Errorf(
		"Config requests user events, but bot does not implement " +
			"UserHandler interface.",
	)

	inboxHandlerErr = fmt.Errorf(
		"Config requests user events, but bot does not implement " +
			"InboxHandler interface.",
	)
)

func subredditStream(
	subs []string,
	loggedIn bool,
	reaper reap.Reaper,
	bot interface{},
) (dispatcher.Dispatcher, error) {
	h, ok := bot.(SubredditHandler)
	if !ok {
		return nil, subredditHandlerErr
	}

	path := subredditsPath(subs)
	if !loggedIn {
		path = logPathsOut([]string{path})[0]
	}

	mon, err := monitor.New(
		monitor.Config{
			Path:   path,
			Lurker: api.NewLurker(reaper),
			Sorter: rsort.New(),
		},
	)

	return dispatcher.New(
		dispatcher.Config{
			Monitor: mon,
			PostHandler: handlers.PostHandlerFunc(
				func(p *data.Post) error {
					return h.Post((*Post)(p))
				},
			),
		},
	), err
}

func userStreams(
	users []string,
	loggedIn bool,
	reaper reap.Reaper,
	bot interface{},
) ([]dispatcher.Dispatcher, error) {
	h, ok := bot.(UserHandler)
	if !ok {
		return nil, userHandlerErr
	}

	paths := userPaths(users)
	if !loggedIn {
		paths = logPathsOut(paths)
	}

	lurker := api.NewLurker(reaper)
	sorter := rsort.New()
	dispatchers := make([]dispatcher.Dispatcher, len(paths))
	for i, p := range paths {
		mon, err := monitor.New(
			monitor.Config{
				Path:   p,
				Lurker: lurker,
				Sorter: sorter,
			},
		)
		if err != nil {
			return nil, err
		}

		dispatchers[i] = dispatcher.New(
			dispatcher.Config{
				Monitor: mon,
				PostHandler: handlers.PostHandlerFunc(
					func(p *data.Post) error {
						return h.UserPost((*Post)(p))
					},
				),
				CommentHandler: handlers.CommentHandlerFunc(
					func(c *data.Comment) error {
						return h.UserComment(
							(*Comment)(c),
						)
					},
				),
			},
		)
	}

	return dispatchers, nil
}

func inboxStream(
	reaper reap.Reaper,
	bot interface{},
) (dispatcher.Dispatcher, error) {
	h, ok := bot.(InboxHandler)
	if !ok {
		return nil, inboxHandlerErr
	}

	mon, err := monitor.New(
		monitor.Config{
			Path:   "/message/inbox",
			Lurker: api.NewLurker(reaper),
			Sorter: rsort.New(),
		},
	)
	if err != nil {
		return nil, err
	}

	router := func(m *data.Message) error {
		if m.Subject == "comment reply" && m.WasComment {
			return h.CommentReply((*Message)(m))
		} else if m.Subject == "post reply" && m.WasComment {
			return h.PostReply((*Message)(m))
		} else if m.Subject == "username mention" && m.WasComment {
			return h.Mention((*Message)(m))
		}

		return h.Message((*Message)(m))
	}

	return dispatcher.New(
		dispatcher.Config{
			Monitor:        mon,
			MessageHandler: handlers.MessageHandlerFunc(router),
		},
	), nil
}
