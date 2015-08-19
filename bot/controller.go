package bot

// Controller defines the interface for bots to interact with the engine. These
// methods are requests to the engine to perform actions on behalf of the bot,
// when it decides it is time.
type Controller interface {
	// Stop stops the engine execution.
	Stop()
}
