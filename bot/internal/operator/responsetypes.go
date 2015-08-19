package operator

import (
	"github.com/turnage/redditproto"
)

// Links returns the links contained in a linkListing.
func getLinks(listing *redditproto.LinkListing) []*redditproto.Link {
	links := make([]*redditproto.Link, len(listing.Data.Children))
	for i, child := range listing.Data.Children {
		links[i] = child.Data
	}
	return links
}
