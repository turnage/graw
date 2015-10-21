package scanner

import (
	"github.com/turnage/redditproto"
)

type Scanner interface {
	Scan() ([]*redditproto.Link, []*redditproto.Comment, error)
}
