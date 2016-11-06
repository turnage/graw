// Package engine runs dispatchers.
package engine

import (
	"time"

	"github.com/turnage/graw/internal/dispatcher"
)

// Config configures an engine.
type Config struct {
	// Dispatchers the engine will control.
	Dispatchers []dispatcher.Dispatcher
	// Rate limits the rate at which dispatchers run.
	Rate <-chan time.Time
}

// Engine controls disptachers.
type Engine interface {
	// Run starts the engine cycle, which runs until it encounters an error
	// or Stop() is called.
	Run() error
	// Stop stops the engine cycle. If the engine is not running, Stop is a
	// no-op.
	Stop()
}

type engine struct {
	ds   []dispatcher.Dispatcher
	rate <-chan time.Time
	stop chan bool
}

// New returns an Engine implementation.
func New(c Config) Engine {
	return &engine{
		ds:   c.Dispatchers,
		rate: c.Rate,
		stop: make(chan bool, 100),
	}
}

func (e *engine) Run() error {
	var dispatcher int = 0
	i := func() int {
		defer func() {
			dispatcher++
			dispatcher %= len(e.ds)
		}()
		return dispatcher
	}
	for {
		select {
		case <-e.rate:
			if err := e.ds[i()].Dispatch(); err != nil {
				return err
			}
		case <-e.stop:
			return nil
		}
	}
}

func (e *engine) Stop() {
	e.stop <- true
}
