// Package api defines the graw api for Reddit bots.
package api

// Actor defines methods for bots that do things (send messages, make posts,
// fetch threads, etc).
type Actor interface {
	// TakeEngine is called when the engine starts; bots should save the
	// engine so they can call its methods. This is only called once.
	TakeEngine(eng Engine)
}

// Loader defines methods for bots that use external resources or need to do
// initialization.
type Loader interface {
	// SetUp is the first method ever called on the bot, and it will be
	// allowed to finish before other methods are called. Bots should
	// load resources here.
	SetUp() error
	// TearDown is the last method ever called on the bot, and all other
	// method calls will finish before this method is called. Bots should
	// unload resources here.
	TearDown() error
}

// Failer defines methods bots can use to control how the Engine responds to
// failures.
type Failer interface {
	// Fail will be called when the engine encounters an error. The bot can
	// return true to instruct the engine to fail, or false to instruct the
	// engine to try again.
	//
	// This method will be called in the main engine loop; the bot may
	// choose to pause here or do other things to respond to the failure
	// (e.g. pause for three hours to respond to Reddit down time).
	Fail(err error) bool
}
