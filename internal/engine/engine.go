// Package engine runs dispatchers.
package engine

import (
	"log"
	"time"

	"github.com/turnage/graw/internal/client"
	"github.com/turnage/graw/internal/dispatcher"
)

// Config configures an engine.
type Config struct {
	// Dispatchers the engine will control.
	Dispatchers []dispatcher.Dispatcher
	// Rate limits the rate at which dispatchers run.
	Rate <-chan time.Time
	// Logger logs events and errors.
	Logger *log.Logger
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
	logger *log.Logger
	ds     []dispatcher.Dispatcher
	rate   <-chan time.Time
	stop   chan bool
}

// New returns an Engine implementation.
func New(c Config) Engine {
	return &engine{
		logger: c.Logger,
		ds:     c.Dispatchers,
		rate:   c.Rate,
		stop:   make(chan bool, 100),
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
			if len(e.ds) == 0 {
				break
			}
			err := e.ds[i()].Dispatch()
			switch err {
			case client.BusyErr:
				e.logger.Printf("503: Busy from Reddit; ignoring")
			case client.GatewayErr:
				e.logger.Printf("502: Bad Gateway from Reddit; ignoring")
			case nil:
			default:
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
