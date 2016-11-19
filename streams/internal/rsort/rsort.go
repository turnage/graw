//go:generate genny -in=thingcast.tpl -out=postcast.go gen "ThingType=*reddit.Post name=posts NAME=Post"
//go:generate genny -in=thingcast.tpl -out=commentcast.go gen "ThingType=*reddit.Comment name=comments NAME=Comment"
//go:generate genny -in=thingcast.tpl -out=messagecast.go gen "ThingType=*reddit.Message name=messages NAME=Message"
// Package rsort provides tools for sorting Reddit elements.
package rsort

import (
	"sort"

	"github.com/turnage/graw/reddit"
)

// Sorter sorts Reddit element harvests.
type Sorter interface {
	// Sort sorts a Reddit element harvest and returns its fullnames in the
	// order of their creation (younger names first).
	Sort(h reddit.Harvest) []string
}

type sorter struct{}

func (s *sorter) Sort(h reddit.Harvest) []string {
	return sortHarvest(h)
}

// New returns a new sorter implementation.
func New() Sorter {
	return &sorter{}
}

// redditThing is named after the Reddit class "Thing", from which all items
// with a full name and creation time inherit.
type redditThing interface {
	Name() string
	Birth() uint64
}

// sortHarvest returns the list of names of Reddit elements in a harvest sorted
// by creation time to the younger elements appear first in the slice.
func sortHarvest(h reddit.Harvest) []string {
	things := merge(
		postsAsThings(h.Posts),
		commentsAsThings(h.Comments),
		messagesAsThings(h.Messages),
	)
	sort.Sort(byCreationTime{things})

	names := make([]string, len(things))
	for i, t := range things {
		names[i] = t.Name()
	}

	return names
}

func merge(things ...[]redditThing) []redditThing {
	var result []redditThing
	for _, t := range things {
		result = append(result, t...)
	}
	return result
}
