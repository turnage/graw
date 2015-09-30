// Package scanner provides Scanner, which scans listing endpoints on reddit.com
// for new Things (instances of the unfortunately named Reddit class, Thing).
package scanner

import (
	"github.com/turnage/graw/internal/operator"
)

const (
	// maxTipSize is the number of posts to keep in the tracked tip. > 1
	// is kept because a tip is needed to fetch only posts newer than
	// that post. If one is deleted, PostMonitor moves to a fallback tip.
	maxTipSize = 15
	// defaultBlankThreshold is the number of refreshes returning no
	// listings the scanner will forgive before attempting to fix its tip.
	defaultBlankThreshold = 1
)

// Scanner scans a listing endpoint on reddit.com for new Things.
type Scanner struct {
	// op is the operator Scanner will make requests to reddit with.
	op operator.Operator
	// blanks is the number of consecutive queries which have returned no
	// Things.
	blanks int
	// blankThreshold is the number of consevutive queries returning no
	// Things which will be tolerated before the Scanner distrusts its
	// current tip; this is adjusted for the Scanner because some paths are
	// less active than others.
	blankThreshold int
	// path is the url path of a listing on reddit.com this Scanner will
	// monitor.
	path string
	// kind is the kind of listing path accesses.
	kind operator.Kind
	// tip is the list of most recent Thing fullnames from the listing this
	// Scanner monitors. This is used because reddit allows queries relative
	// to fullnames.
	tip []string
}

// New returns a Scanner which scans the given listing path using the given
// operator.
func New(path string, op operator.Operator, kind operator.Kind) *Scanner {
	return &Scanner{
		op:             op,
		blankThreshold: defaultBlankThreshold,
		path:           path,
		kind:           kind,
		tip:            []string{""},
	}
}

// Scan returns new elements in the listing.
func (s *Scanner) Scan() ([]operator.Thing, error) {
	things, err := s.fetchTip()
	if err != nil {
		return nil, err
	}

	if len(things) == 0 {
		s.blanks++
		if s.blanks > s.blankThreshold {
			shaved, err := s.fixTip()
			if err != nil {
				return nil, err
			}
			if !shaved {
				s.blankThreshold += defaultBlankThreshold
			}
			s.blanks = 0
		}
	} else {
		s.blanks = 0
	}

	return things, nil
}

// fetchTip fetches the latest posts from the monitored subreddits. If there is
// no tip, fetchTip considers the call an adjustment round, and will fetch a new
// reference tip but discard the post (because, most likely, that post was
// already returned before).
func (s *Scanner) fetchTip() ([]operator.Thing, error) {
	tip := s.tip[len(s.tip)-1]
	links := uint(operator.MaxLinks)
	adjustment := false
	if tip == "" {
		links = 1
		adjustment = true
	}

	things, err := s.op.Scrape(
		s.path,
		"",
		tip,
		links,
		operator.Link,
	)
	if err != nil {
		return nil, err
	}

	for i := range things {
		s.tip = append(s.tip, things[len(things)-1-i].GetName())
	}

	if len(s.tip) > maxTipSize {
		s.tip = s.tip[len(s.tip)-maxTipSize:]
	}

	if adjustment && len(things) == 1 {
		return nil, nil
	}

	return things, nil
}

// fixTip attempts to fix the PostMonitor's reference point for new posts. If it
// has been deleted, fixTip will move to a fallback tip. fixtip returns true if
// the tip was shaved.
func (s *Scanner) fixTip() (bool, error) {
	thing, err := s.op.GetThing(s.tip[len(s.tip)-1], s.kind)
	if err != nil {
		return false, err
	}

	if thing == nil {
		s.shaveTip()
		return true, nil
	}

	return false, nil
}

// shaveTip shaves off the latest tip thread name. If all tips are shaved off,
// uses an empty tip name (this will just get the latest threads).
func (s *Scanner) shaveTip() {
	if len(s.tip) == 1 {
		s.tip[0] = ""
		return
	}

	s.tip = s.tip[:len(s.tip)-1]
}
