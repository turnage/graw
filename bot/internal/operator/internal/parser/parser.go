// Package parser parses Reddit's JSON responses into protobuffers.
package parser

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/turnage/redditproto"
)

// ParseLinkListing returns a slice of Links which hold the same data the JSON
// link listing provided contains.
func ParseLinkListing(content io.ReadCloser) ([]*redditproto.Link, error) {
	if content == nil {
		return nil, fmt.Errorf("no content provided")
	}

	listing := &redditproto.LinkListing{}
	decoder := json.NewDecoder(content)
	if err := decoder.Decode(listing); err != nil {
		return nil, err
	}

	return unpackLinkListing(listing)
}

// ParseThread parses a combination link listing and comment listing, which
// Reddit returns when asked for the JSON digest of a thread. This contains the
// submission's information, and all of its comments. The returned link will
// have the Comments field filled, and the comments will have their ReplyTree
// field filled.
func ParseThread(content io.ReadCloser) (*redditproto.Link, error) {
	if content == nil {
		return nil, fmt.Errorf("no content provided")
	}

	listings := []interface{}{
		&redditproto.LinkListing{},
		&redditproto.CommentListing{},
	}
	decoder := json.NewDecoder(content)
	if err := decoder.Decode(&listings); err != nil {
		return nil, err
	}

	if len(listings) != 2 {
		return nil, fmt.Errorf(
			"json decoding malformed the listings: %v",
			listings)
	}

	linkListing := listings[0].(*redditproto.LinkListing)
	commentListing := listings[1].(*redditproto.CommentListing)

	unpackedLinks, err := unpackLinkListing(linkListing)
	if err != nil {
		return nil, err
	}

	if len(unpackedLinks) != 1 {
		return nil, fmt.Errorf(
			"unexpected amount of links (%d)",
			len(unpackedLinks))
	}

	link := unpackedLinks[0]
	link.Comments = unpackCommentListing(commentListing)

	return link, nil
}

// unpackLinkListing returns a slice of the links contained in a link listing.
func unpackLinkListing(
	listing *redditproto.LinkListing,
) ([]*redditproto.Link, error) {
	if listing.GetData() == nil {
		return nil, fmt.Errorf("no data field; got %v", listing)
	}

	if listing.GetData().GetChildren() == nil {
		return nil, fmt.Errorf("data has no children; got %v", listing)
	}

	links := make([]*redditproto.Link, len(listing.GetData().GetChildren()))
	for i, child := range listing.GetData().GetChildren() {
		links[i] = child.GetData()
	}
	return links, nil
}

// unpackCommentListing returns a slice of the comments contained in a comment
// listing.
func unpackCommentListing(
	listing *redditproto.CommentListing,
) []*redditproto.Comment {
	if listing.GetData() == nil {
		return nil
	}

	if listing.GetData().GetChildren() == nil {
		return nil
	}

	comments := make(
		[]*redditproto.Comment,
		len(listing.GetData().GetChildren()))
	for i, child := range listing.GetData().GetChildren() {
		comments[i] = child.GetData()
		if comments[i].Replies != nil {
			comments[i].ReplyTree = unpackCommentListing(comments[i].Replies)
			comments[i].Replies = nil
		}
	}
	return comments
}
