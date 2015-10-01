package operator

import (
	"github.com/turnage/redditproto"
)

// MockOperator mocks Operator; it returns canned responses.
type MockOperator struct {
	// ScrapeErr is returned in the error field of Scrape.
	ScrapeErr error
	// ScrapeReturn is returned by Scrape.
	ScrapeReturn []Thing
	// GetThingErr is returned in the error field of GetThing.
	GetThingErr error
	// GetThingReturn is returned by GetThing.
	GetThingReturn Thing
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
	kind Kind,
) ([]Thing, error) {
	return m.ScrapeReturn, m.ScrapeErr
}

func (m *MockOperator) GetThing(id string, kind Kind) (Thing, error) {
	return m.GetThingReturn, m.GetThingErr
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
