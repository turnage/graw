package graw

// Engine defines the interface for bots to interact with the engine. These
// methods are requests to the engine to perform actions on behalf of the bot,
// when it decides it is time.
type Engine interface {
	// Stop stops the engine execution.
	Stop()
}
