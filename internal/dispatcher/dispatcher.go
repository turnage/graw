// Package dispatcher connects a monitor to handlers which process the new
// events.
package dispatcher

import (
	"golang.org/x/sync/errgroup"

	"github.com/turnage/graw/internal/handlers"
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/reap"
)

// Dispatcher feeds new events from its monitor to its handlers.
type Dispatcher interface {
	// Dispatch pulls new events from the Dispatcher's monitor and forwards
	// them to their corresponding handlers.
	Dispatch() error
}

// Config configures a dispatcher.
type Config struct {
	// Monitor is the source of events.
	Monitor monitor.Monitor

	// Handlers process events from the harvest.

	PostHandler    handlers.PostHandler
	CommentHandler handlers.CommentHandler
	MessageHandler handlers.MessageHandler
}

type dispatcher struct {
	mon monitor.Monitor
	ph  handlers.PostHandler
	ch  handlers.CommentHandler
	mh  handlers.MessageHandler
}

func New(c Config) Dispatcher {
	return &dispatcher{
		mon: c.Monitor,
		ph:  c.PostHandler,
		ch:  c.CommentHandler,
		mh:  c.MessageHandler,
	}
}

func (d *dispatcher) Dispatch() error {
	harvest, err := d.mon.Update()
	if err != nil {
		return err
	}

	return d.dispatch(harvest)
}

func (d *dispatcher) dispatch(h reap.Harvest) error {
	var g errgroup.Group

	// LOL NO GENERICS

	if d.ph != nil {
		for i := range h.Posts {
			p := h.Posts[i]
			wrap := func() error { return d.ph.HandlePost(p) }
			g.Go(wrap)
		}
	}

	if d.ch != nil {
		for i := range h.Comments {
			c := h.Comments[i]
			wrap := func() error { return d.ch.HandleComment(c) }
			g.Go(wrap)
		}
	}

	if d.mh != nil {
		for i := range h.Messages {
			m := h.Messages[i]
			wrap := func() error { return d.mh.HandleMessage(m) }
			g.Go(wrap)
		}
	}

	return g.Wait()
}
