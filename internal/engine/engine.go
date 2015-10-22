// Package engine provides implementations for bot engines. See the provider
// functions for details about what context they run a bot in.
package engine

// Ignition provides an interface to start the engine.
type Ignition interface {
	// Run should be called once to start the engine. It may run forever.
	Run() error
}
