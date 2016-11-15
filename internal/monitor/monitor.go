// Package monitor includes monitors for different parts of Reddit, such as a
// user inbox or a subreddit's post feed.
package monitor

import (
	"github.com/turnage/graw/internal/api/scanner"
	"github.com/turnage/graw/internal/reap"
	"github.com/turnage/graw/internal/rsort"
)

const (
	// The blank threshold is the amount of updates returning 0 new
	// elements in the monitored listing the monitor will tolerate before
	// suspecting the tip of the listing has been deleted.
	blankThreshold = 1
	// maxTipSize is the maximum size of the tip log (number of backup tips
	// + the current tip).
	maxTipSize = 20
)

// Monitor defines the controls for a Monitor.
type Monitor interface {
	// Update will check for new events, and send them to the Monitor's
	// handlers.
	Update() (reap.Harvest, error)
}

// Config configures a monitor.
type Config struct {
	// Path is the path to the listing the monitor watches.
	Path string

	// Scanner is the api the monitor uses to read Reddit data.
	Scanner scanner.Scanner

	// Sorter sorts the monitor's new listing elements.
	Sorter rsort.Sorter
}

type monitor struct {
	// blanks is the number of rounds that have turned up 0 new
	// elements at the listing endpoint.
	blanks int
	// blankThreshold is the number of blanks a monitor will tolerate before
	// suspecting its tip is broken (e.g.post was deleted).
	blankThreshold int
	// tip is a slice of reddit thing names, the first of which represents
	// the "tip", which the monitor uses to requests new posts by using it
	// as a reference point (i.e.asks Reddit for posts "after" the tip).
	tip []string
	// path is the listing endpoint the monitor monitors. This path is
	// appended to the reddit monitor url (e.g./user/robert).
	path string

	scanner scanner.Scanner
	sorter  rsort.Sorter
}

// New provides a monitor for the listing endpoint.
func New(c Config) (Monitor, error) {
	m := &monitor{
		blankThreshold: blankThreshold,
		tip:            []string{""},
		path:           c.Path,
		scanner:        c.Scanner,
		sorter:         c.Sorter,
	}

	if err := m.sync(); err != nil {
		return nil, err
	}

	return m, nil
}

// Update checks for new content at the monitored listing endpoint and forwards
// new content to the bot for processing.
func (m *monitor) Update() (reap.Harvest, error) {
	if m.blanks == m.blankThreshold {
		return reap.Harvest{}, m.healthCheck()
	}

	harvest, err := m.scanner.Listing(m.path, m.tip[0])
	m.updateTip(harvest)
	return harvest, err
}

// sync fetches the current tip of a listing endpoint, so that grawbots crawling
// forward in time don't treat it as a new post, or reprocess it when restarted.
func (m *monitor) sync() error {
	harvest, err := m.scanner.Listing(m.path, "")
	if err != nil {
		return err
	}

	names := m.sorter.Sort(harvest)
	if len(names) > 0 {
		m.tip = names
	} else {
		m.tip = []string{""}
	}

	return nil
}

// updateTip updates the monitor's list of names from the endpoint listing it
// uses to keep track of its position in the monitored listing (e.g. a user's
// page or its position in a subreddit's history).
func (m *monitor) updateTip(h reap.Harvest) {
	names := m.sorter.Sort(h)
	if len(names) > 0 {
		m.blanks = 0
	} else {
		m.blanks++
	}

	m.tip = append(names, m.tip...)
	if len(m.tip) > maxTipSize {
		m.tip = m.tip[0:maxTipSize]
	}
}

// healthCheck checks the health of the tip when nothing is returned from a
// scrape enough times.
func (m *monitor) healthCheck() error {
	m.blanks = 0
	broken, err := m.fixTip()
	if err != nil {
		return err
	}
	if !broken {
		m.blankThreshold *= 2
	}

	return nil
}

// fixTip checks that the fullname at the front of the tip is still valid (e.g.
// not deleted).If it isn't, it shaves the tip.fixTip returns whether the tip
// was broken.
func (m *monitor) fixTip() (bool, error) {
	exists, err := m.scanner.Exists(m.tip[0])
	if err != nil {
		return false, err
	}

	if !exists {
		m.shaveTip()
	}

	return !exists, nil
}

// shaveTip shaves the latest fullname off of the tip, promoting the preceding
// fullname if there is one or resetting the tip if there isn't.
func (m *monitor) shaveTip() {
	if len(m.tip) <= 1 {
		m.tip = []string{""}
	} else {
		m.tip = m.tip[1:]
	}
}
