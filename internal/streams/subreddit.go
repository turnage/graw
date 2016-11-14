package streams

import (
	"fmt"

	"github.com/turnage/graw/internal/handlers"
)

func noHandlerErr(eventType string) error {
	return fmt.Errorf(
		"%s events requested, but no handler provided.",
		eventType,
	)
}

func subreddits(subs []string, sh handlers.SubredditHandler) (
	singleConfig,
	error,
) {
	if sh == nil {
		return singleConfig{}, noHandlerErr("Subreddit")
	}

	return singleConfig{
		path: subredditsPath(subs),
		ph:   handlers.DecomposeSubredditHandler(sh),
	}, nil
}

func users(users []string, uh handlers.UserHandler) ([]singleConfig, error) {
	if uh == nil {
		return nil, noHandlerErr("User")
	}

	paths := userPaths(users)
	ph, ch := handlers.DecomposeUserHandler(uh)
	configs := []singleConfig{}
	for _, p := range paths {
		configs = append(
			configs, singleConfig{
				path: p,
				ph:   ph,
				ch:   ch,
			},
		)
	}

	return configs, nil
}

func inbox(loggedIn bool, ih handlers.InboxHandler) (singleConfig, error) {
	if ih == nil {
		return singleConfig{}, noHandlerErr("Inbox")
	}

	if !loggedIn {
		return singleConfig{}, fmt.Errorf(
			"Inbox events requested, but the bot is not logged in!",
		)
	}

	return singleConfig{
		path: "/message/inbox",
		mh:   handlers.DecomposeInboxHandler(ih),
	}, nil
}
