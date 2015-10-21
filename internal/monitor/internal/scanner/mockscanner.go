package scanner

import (
	"github.com/turnage/redditproto"
)

// MockScanner provides canned responses to all calls to Scanner methods.
type MockScanner struct {
	ScanLinksReturn    []*redditproto.Link
	ScanCommentsReturn []*redditproto.Comment
	ScanErr            error
}

// Scan returns canned responses to the method call.
func (m *MockScanner) Scan() (
	[]*redditproto.Link,
	[]*redditproto.Comment,
	error,
) {
	return m.ScanLinksReturn, m.ScanCommentsReturn, m.ScanErr
}
