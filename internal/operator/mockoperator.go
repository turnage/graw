package operator

import (
	"github.com/turnage/redditproto"
)

// MockOperator mocks Operator; it returns canned responses.
type MockOperator struct {
	Err           error
	ScrapeReturn  []*redditproto.Link
	ThreadsReturn []*redditproto.Link
	ThreadReturn  *redditproto.Link
	InboxReturn   []*redditproto.Message
}

func (m *MockOperator) Scrape(
	subreddit,
	sort,
	after,
	before string,
	limit uint,
) ([]*redditproto.Link, error) {
	return m.ScrapeReturn, m.Err
}

func (m *MockOperator) Threads(
	fullnames ...string,
) ([]*redditproto.Link, error) {
	return m.ThreadsReturn, m.Err
}

func (m *MockOperator) Thread(permalink string) (*redditproto.Link, error) {
	return m.ThreadReturn, m.Err
}

func (m *MockOperator) Inbox() ([]*redditproto.Message, error) {
	return m.InboxReturn, m.Err
}

func (m *MockOperator) Reply(parent, content string) error {
	return m.Err
}

func (m *MockOperator) Submit(subreddit, kind, title, content string) error {
	return m.Err
}

func (m *MockOperator) Compose(user, subject, content string) error {
	return m.Err
}
