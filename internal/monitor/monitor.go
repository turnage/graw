// Package monitor continuously updates monitored sections of reddit, such as
// subreddits and threads.
package monitor

import (
	"bytes"
	"container/list"
	"sync"
	"time"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

const (
	// errorTolerance is the number of errors the monitor will tolerate from
	// reddit in a row before it shuts down. Monitor backs off exponentially
	// between the allowed errors.
	errorTolerance = 5
	// maxPosts is the maximum number of posts to request new at once.
	maxPosts = 100
	// maxTipSize is the number of posts to keep in the tracked tip. More than
	// one is kept because a tip is needed to fetch only posts newer than
	// that post. If one is deleted, monitor moves to a fallback tip.
	maxTipSize = 15
)

// Monitor monitors sections of reddit real time and exports updates. All
// methods and channels of Monitor expect that Run() is alive in a goroutine.
// Calling them when that condition is not true is not defined behavior.
type Monitor struct {
	// NewPosts provides new posts to monitored subreddits. These posts will
	// have been posted very recently so they probably won't have comments
	// or votes yet.
	NewPosts chan *redditproto.Link
	// NewMessages provides new private messages to the bot's inbox.
	NewMessages chan *redditproto.Message
	// NewCommentReplies provides new comment replies to the bot's inbox.
	NewCommentReplies chan *redditproto.Message
	// NewPostReplies provides new post replies to the bot's inbox.
	NewPostReplies chan *redditproto.Message
	// NewMentions provides new mentions of the bot's username.
	NewMentions chan *redditproto.Message
	// Errors provides errors that cause monitor to quit running.
	Errors chan error

	// op is the operator through which the monitor will make update
	// requests to reddit.
	op *operator.Operator
	// tip is the list of latest posts in the monitored subreddits.
	tip *list.List
	// errors is a count of how many errors monitor has encountered trying
	// to talk to reddit.
	errors uint
	// errorBackoffUnit is the unit of time that the error back off strategy
	// will increase and sleep between consecutive errors.
	errorBackOffUnit time.Duration
	// blanks is a count of how many times monitor has updated and found no
	// new posts in row.
	blanks uint
	// blankRoundTolerance is how many times no new posts can be found in
	// monitored subreddits before monitor attempts to fix its tip.
	blankRoundTolerance uint

	// mu protects the following fields.
	mu sync.Mutex
	// monitoredSubreddits is the list of monitored subreddits from which
	// the requests are built.
	monitoredSubreddits map[string]bool
	// subredditQuery is the subreddit alias Monitor uses to fetch new
	// posts. It uses reddit's "+" multireddit technique, e.g. "self+aww".
	subredditQuery string
}

// New returns an initialized Monitor.
func New(op *operator.Operator, subreddits []string) *Monitor {
	mon := &Monitor{
		NewPosts:            make(chan *redditproto.Link),
		NewMessages:         make(chan *redditproto.Message),
		NewMentions:         make(chan *redditproto.Message),
		NewCommentReplies:   make(chan *redditproto.Message),
		NewPostReplies:      make(chan *redditproto.Message),
		Errors:              make(chan error),
		errorBackOffUnit:    time.Minute,
		op:                  op,
		tip:                 list.New(),
		monitoredSubreddits: make(map[string]bool),
	}
	mon.tip.PushFront("")
	mon.MonitorSubreddits(subreddits...)
	return mon
}

// Run is the main loop of the monitor, and output is fed through
// Monitor's exported channels.
func (m *Monitor) Run() {
	for true {
		postCount, err := m.updatePosts()
		if m.errorBackOff(err) {
			return
		}
		m.checkOnTip(postCount)
		if m.errorBackOff(m.updateInbox()) {
			return
		}
	}
}

// MonitorSubreddits starts monitoring the requested subreddits.
func (m *Monitor) MonitorSubreddits(subreddits ...string) {
	m.mu.Lock()
	setKeys(m.monitoredSubreddits, true, subreddits)
	m.subredditQuery = buildQuery(m.monitoredSubreddits, "+")
	m.mu.Unlock()
}

// UnmonitorSubreddits stops monitoring the requested subreddits.
func (m *Monitor) UnmonitorSubreddits(subreddits ...string) {
	m.mu.Lock()
	setKeys(m.monitoredSubreddits, false, subreddits)
	m.subredditQuery = buildQuery(m.monitoredSubreddits, "+")
	m.mu.Unlock()
}

// checkOnTip keeps track of how many times no posts have been returned on a
// scrape, and once that has exceeded the tolerance, attempts to fix the tip.
func (m *Monitor) checkOnTip(postCount int) error {
	if postCount > 0 {
		return nil
	}

	m.blanks++
	if m.blanks > m.blankRoundTolerance {
		broken, err := m.fixTip()
		if err != nil {
			return err
		}
		if !broken {
			m.blankRoundTolerance++
		}
		m.blanks = 0
	}

	return nil
}

// errorBackOff keeps a count of errors, and backs off by blocking the monitor
// thread based on how many errors have occurred consecutively.
//
// It returns whether the error tolerance has been exceeded.
func (m *Monitor) errorBackOff(err error) bool {
	if err == nil {
		m.errors = 0
		return false
	}

	m.errors++
	if m.errors > errorTolerance {
		m.Errors <- err
		return true
	}

	time.Sleep(m.errorBackOffUnit << m.errors)
	return false
}

// updateSubreddits gets new posts from monitored subreddits and feeds them over
// the output channel.
func (m *Monitor) updatePosts() (int, error) {
	posts, err := m.fetchTip()
	if err != nil {
		return 0, err
	}
	for _, post := range posts {
		m.NewPosts <- post
	}
	return len(posts), nil
}

// updateInbox gets unread messages from monitored subreddits and feeds them
// over the output channel.
func (m *Monitor) updateInbox() error {
	messages, err := m.op.Inbox()
	if err != nil {
		return err
	}
	for _, message := range messages {
		if message.GetSubject() == "username mention" {
			m.NewMentions <- message
		} else if message.GetSubject() == "post reply" {
			m.NewPostReplies <- message
		} else if message.GetWasComment() {
			m.NewCommentReplies <- message
		} else {
			m.NewMessages <- message
		}
	}
	return nil
}

// fetchTip fetches the latest posts from the monitored subreddits.
func (m *Monitor) fetchTip() ([]*redditproto.Link, error) {
	posts, err := m.op.Scrape(
		m.subredditQuery,
		"new",
		"",
		m.tip.Front().Value.(string),
		maxPosts,
	)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		m.tip.PushFront(posts[len(posts)-1-i].GetName())
		if m.tip.Len() > maxTipSize {
			m.tip.Remove(m.tip.Back())
		}
	}

	return posts, nil
}

// fixTip fixes the tip if the post has been deleted. fixTip returns whether
// the tip was broken.
func (m *Monitor) fixTip() (bool, error) {
	wasBroken := false
	ids := make([]string, m.tip.Len())
	for e := m.tip.Front(); e != nil; e = e.Next() {
		ids = append(ids, e.Value.(string))
	}
	posts, err := m.op.Threads(ids...)
	if err != nil {
		return false, err
	}

	for e := m.tip.Front(); e != nil; e = e.Next() {
		if e.Prev() != nil {
			wasBroken = true
			m.tip.Remove(e.Prev())
		}
		for _, post := range posts {
			if e.Value.(string) == post.GetName() {
				return wasBroken, nil
			}
		}
	}
	m.tip.Remove(m.tip.Front())
	m.tip.PushFront("")

	return wasBroken, nil
}

// setKeys sets the value of all provided keys to val in m.
func setKeys(m map[string]bool, val bool, keys []string) {
	for _, key := range keys {
		m[key] = val
	}
}

// buildQuery assembles a delimited list of some kind of name to use as a query
// to reddit, from a map indicating whether each name should be included in the
// query.
func buildQuery(names map[string]bool, delim string) string {
	var queryBuffer bytes.Buffer
	emptyQuery := true

	for name, include := range names {
		if include {
			emptyQuery = false
			queryBuffer.WriteString(name)
			queryBuffer.WriteString(delim)
		}
	}

	if emptyQuery {
		return ""
	}
	query := queryBuffer.String()

	return query[:len(query)-len(delim)]
}
