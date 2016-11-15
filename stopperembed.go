package graw

import (
	"github.com/turnage/graw/internal/engine"
)

// Stopper can stop the graw engine. Embed it in your bot to access Stop().
type Stopper struct {
	engine engine.Engine
}

type stopper interface {
	grawSetEngine(engine.Engine)
	Stop()
}

func (s *Stopper) grawSetEngine(e engine.Engine) {
	s.engine = e
}

func (s *Stopper) Stop() {
	s.engine.Stop()
}
