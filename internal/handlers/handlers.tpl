package handlers

import (
	"github.com/cheekybits/genny/generic"

	"github.com/turnage/graw/internal/data"
)

type name generic.Type
type NAME generic.Type
type EventType generic.Type

// NAMEHandler processes Reddit names.
type NAMEHandler interface {
	// HandleNAME processes a Reddit name.
	HandleNAME(p EventType) error
}

type nameHandler struct {
	handler func(p EventType) error
}

func (h *nameHandler) HandleNAME(name EventType) error {
	return h.handler(name)
}

// NAMEHandlerFunc returns a NAMEHandler using the given function.
func NAMEHandlerFunc(f func(e EventType) error) NAMEHandler {
	return &nameHandler{f}
}
