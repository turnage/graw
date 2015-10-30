package operator

import (
	"github.com/turnage/redditproto"
)

// MockOperator mocks Operator; it returns canned responses.
type MockOperator struct {
	// ScrapeErr is returned in the error field of Scrape.
	ScrapeErr error
	// ScrapeLinksReturn is returned by Scrape.
	ScrapeLinksReturn []*redditproto.Link
	// ScrapeCommentsReturn is returned by Scrape.
	ScrapeCommentsReturn []*redditproto.Comment
	// ScrapeMessagesReturn is returned by Scrape.
	ScrapeMessagesReturn []*redditproto.Message
	// IsThereThingErr is returned in the error field of IsThereThing.
	IsThereThingErr error
	// IsThereThingReturn is returned by IsThereThing.
	IsThereThingReturn bool
	// ThreadErr is returned in the error field of Thread.
	ThreadErr error
	// ThreadReturn is returned by Thread.
	ThreadReturn *redditproto.Link
	// InboxErr is returned in the error field of Inbox.
	InboxErr error
	// InboxReturn is returned by Inbox.
	InboxReturn []*redditproto.Message
	// MarkAsReadErr is returned in the error field of MarkAsRead.
	MarkAsReadErr error
	// ReplyErr is returned in the error field of Reply.
	ReplyErr error
	// SubmitErr is returned in the error field of Submit.
	SubmitErr error
	// ComposeErr is returned in the error field of Compose.
	ComposeErr error
}

func (m *MockOperator) Scrape(
	path,
	after,
	before string,
	limit uint,
) (
	[]*redditproto.Link,
	[]*redditproto.Comment,
	[]*redditproto.Message,
	error,
) {
	return m.ScrapeLinksReturn,
		m.ScrapeCommentsReturn,
		m.ScrapeMessagesReturn,
		m.ScrapeErr
}

func (m *MockOperator) IsThereThing(id string) (bool, error) {
	return m.IsThereThingReturn, m.IsThereThingErr
}

func (m *MockOperator) Thread(permalink string) (*redditproto.Link, error) {
	return m.ThreadReturn, m.ThreadErr
}

func (m *MockOperator) Inbox() ([]*redditproto.Message, error) {
	return m.InboxReturn, m.InboxErr
}

func (m *MockOperator) MarkAsRead(fullnames ...string) error {
	return m.MarkAsReadErr
}

func (m *MockOperator) Reply(parent, content string) error {
	return m.ReplyErr
}

func (m *MockOperator) Submit(subreddit, kind, title, content string) error {
	return m.SubmitErr
}

func (m *MockOperator) Compose(user, subject, content string) error {
	return m.ComposeErr
}
