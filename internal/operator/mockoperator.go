package operator

import (
	"github.com/turnage/redditproto"
)

// MockOperator mocks Operator; it returns canned responses.
type MockOperator struct {
	ScrapeErr     error
	ScrapeReturn  []Thing
	ThreadsErr    error
	ThreadsReturn []*redditproto.Link
	ThreadErr     error
	ThreadReturn  *redditproto.Link
	InboxErr      error
	InboxReturn   []*redditproto.Message
	MarkAsReadErr error
	ReplyErr      error
	SubmitErr     error
	ComposeErr    error
}

func (m *MockOperator) Scrape(
	path,
	after,
	before string,
	limit uint,
	kind ListingKind,
) ([]Thing, error) {
	return m.ScrapeReturn, m.ScrapeErr
}

func (m *MockOperator) Threads(
	fullnames ...string,
) ([]*redditproto.Link, error) {
	return m.ThreadsReturn, m.ThreadsErr
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
