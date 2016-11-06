package graw

import (
	"time"
)

// App represents an app registered on Reddit in the user preferences "app" tab.
type App struct {
	// ID and Secret are used to claim an OAuth2 grant the users are
	// previously authorized.

	ID     string
	Secret string

	// Username and Password are used to authorize with the endpoint.

	Username string
	Password string
}

// Config is the configuration for a graw run.
type Config struct {
	// Agent is the user agent string to give Reddit. It should be in the
	// form: <platform>:<app ID>:<version string> (by /u/<reddit username>)
	Agent string
	// Subreddits graw should watch for the bot. New posts made in these
	// subreddits will be fed to the bot for processing.
	Subreddits []string
	// Users graw should watch for the bot. New posts and comments made by
	// these users in subreddits the bot has access to will be fed to the
	// bot for processing.
	//
	// Note: Watching many users may make updating other feeds slow, as user
	// feeds must be checked individually and can't be combined in one
	// listing like subreddits.
	Users []string
	// Rate is the interval between request fires to Reddit. The minimum
	// respected value is time.Second.
	Rate time.Duration
	// App is the parameters for identifying as an app registered on Reddit.
	// This is necessary to use the features of a logged in user, but not to
	// use the read-only feeds from Reddit as a logged out user.
	App *App
}
