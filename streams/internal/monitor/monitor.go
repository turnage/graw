// Package monitor tracks a listing feed on Reddit.
package monitor

import (
	"github.com/turnage/graw/reddit"

	"github.com/turnage/graw/streams/internal/rsort"
)

const (
	// The blank threshold is the amount of updates returning 0 new
	// elements in the monitored listing the monitor will tolerate before
	// suspecting the tip of the listing has been deleted or caught in a
	// spam filter.
	blankThreshold = 4
	// maxTipSize is the maximum size of the tip log (number of backup tips
	// + the current tip).
	maxTipSize = 20
)

// defaultTip is the blank reference point in a Reddit listing, which asks for
// the most recent elements.
var defaultTip = []string{""}

// Monitor defines the controls for a Monitor.
type Monitor interface {
	// Update will check for new events, and send them to the Monitor's
	// handlers.
	Update() (reddit.Harvest, error)
}

// Config configures a monitor.
type Config struct {
	// Path is the path to the listing the monitor watches.
	Path string

	// Scanner is the api the monitor uses to read Reddit
	Scanner reddit.Scanner

	// Sorter sorts the monitor's new listing elements.
	Sorter rsort.Sorter
}

type monitor struct {
	// blanks is the number of rounds that have turned up 0 new
	// elements at the listing endpoint.
	blanks int
	// tip is a slice of reddit thing names, the first of which represents
	// the "tip", which the monitor uses to requests new posts by using it
	// as a reference point (i.e.asks Reddit for posts "after" the tip).
	tip []string
	// path is the listing endpoint the monitor monitors. This path is
	// appended to the reddit monitor url (e.g./user/robert).
	path string

	scanner reddit.Scanner
	sorter  rsort.Sorter
}

// New provides a monitor for the listing endpoint.
func New(c Config) (Monitor, error) {
	m := &monitor{
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
func (m *monitor) Update() (reddit.Harvest, error) {
	if m.blanks > blankThreshold {
		return reddit.Harvest{}, m.fixTip()
	}

	names, harvest, err := m.harvest(m.tip[0])
	m.updateTip(names)
	return harvest, err
}

// harvest fetches from the listing any posts after the given reference post,
// and returns those posts and a reverse chronologically sorted list of their
// names.
func (m *monitor) harvest(ref string) ([]string, reddit.Harvest, error) {
	h, err := m.scanner.Listing(m.path, ref)
	return m.sorter.Sort(h), h, err
}

// sync fetches the current tip of a listing endpoint, so that grawbots crawling
// forward in time don't treat it as a new post, or reprocess it when restarted.
func (m *monitor) sync() error {
	names, _, err := m.harvest("")
	if len(names) > 0 {
		m.tip = names
	} else {
		m.tip = defaultTip
	}
	return err
}

// updateTip updates the monitor's list of names from the endpoint listing it
// uses to keep track of its position in the monitored listing.
func (m *monitor) updateTip(names []string) {
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

// fixTip checks all of the stored backup tips for health. If the post at the
// front has been deleted or caught in a spam filter, the feed will die and we
// will stop getting posts. This will adjust backward if a tip is dead and
// remove any other dead tips in the list. Returns whether the tip was broken.
func (m *monitor) fixTip() error {
	names, _, err := m.harvest(m.tip[len(m.tip)-1])
	if err != nil {
		return err
	}

	// If none of our backup tips were returned, most likely the last backup
	// tip is dead and this check was meaningless.
	if len(names) == 0 {
		m.tip = m.tip[:len(m.tip)-1]
		if len(m.tip) == 0 {
			m.tip = defaultTip
		}
		return nil
	}

	// n^2 because your cycles don't matter to me & n <= maxTipSize
	for i := 0; i < len(m.tip)-1; i++ {
		alive := false
		for _, n := range names {
			if m.tip[i] == n {
				alive = true
			}
		}
		if !alive {
			m.tip = append(m.tip[:i], m.tip[i+1:]...)
		}
	}

	m.blanks = 0
	return nil
}
