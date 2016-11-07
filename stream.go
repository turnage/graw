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
