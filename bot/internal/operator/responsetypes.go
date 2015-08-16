package operator

import (
	"github.com/turnage/redditproto"
)

// linkListing is structured in the way Reddit returns listing of links, so that
// they can be unmarshaled into instances of it.
type linkListing struct {
	Data struct {
		Children []struct {
			Data *redditproto.Link
		}
	}
}

// Links returns the links contained in a linkListing.
func (l *linkListing) Links() []*redditproto.Link {
	links := make([]*redditproto.Link, len(l.Data.Children))
	for i, child := range l.Data.Children {
		links[i] = child.Data
	}
	return links
}
