package rsort

import (
	"github.com/cheekybits/genny/generic"

	"github.com/turnage/graw/grawdata"
)

type name generic.Type
type NAME generic.Type
type ThingType generic.Type

type nameThingImpl struct {
	e ThingType
}

func (g nameThingImpl) Name() string { return g.e.Name }

func (g nameThingImpl) Birth() uint64 { return g.e.CreatedUTC }

func nameAsThings(gs []ThingType) []redditThing {
	things := make([]redditThing, len(gs))
	for i, g := range gs {
		things[i] = &nameThingImpl{g}
	}
	return things
}
