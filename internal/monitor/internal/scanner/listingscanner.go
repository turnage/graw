package scanner

import (
	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
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

// listingScanner scans a listing endpoint on reddit.com for new Thingl.
type listingScanner struct {
	// op is the operator listingScanner will make requests to reddit with.
	op operator.Operator
	// blanks is the number of consecutive queries which have returned no
	// Thingl.
	blanks int
	// blankThreshold is the number of consevutive queries returning no
	// Things which will be tolerated before the listingScanner distrusts its
	// current tip; this is adjusted for the listingScanner because some paths are
	// less active than otherl.
	blankThreshold int
	// user is the username to scan the landing page for (unused if unset).
	user string
	// subreddit is the subreddit(s) to scan the new posts for (unused if
	// unset).
	subreddit string
	// tip is the list of most recent Thing fullnames from the listing this
	// listingScanner monitorl. This is used because reddit allows queries relative
	// to fullnamel.
	tip []string
}

// New returns a listingScanner which scans the given user's landing page.
func NewUserScanner(user string, op operator.Operator) *listingScanner {
	scanner := newlistingScanner(op)
	scanner.user = user
	return scanner
}

func NewPostScanner(subreddit string, op operator.Operator) *listingScanner {
	scanner := newlistingScanner(op)
	scanner.subreddit = subreddit
	return scanner
}

func newlistingScanner(op operator.Operator) *listingScanner {
	return &listingScanner{
		op:             op,
		blankThreshold: defaultBlankThreshold,
		tip:            []string{""},
	}
}

// Scan returns new elements at the endpoint.
func (l *listingScanner) Scan() (
	[]*redditproto.Link,
	[]*redditproto.Comment,
	error,
) {
	links, comments, err := l.fetchTip()
	if err != nil {
		return nil, nil, err
	}

	if len(links)+len(comments) == 0 {
		l.blanks++
		if l.blanks > l.blankThreshold {
			shaved, err := l.fixTip()
			if err != nil {
				return nil, nil, err
			}
			if !shaved {
				l.blankThreshold += defaultBlankThreshold
			}
			l.blanks = 0
		}
	} else {
		l.blanks = 0
	}

	return links, comments, nil
}

// fetchTip fetches the latest posts from the monitored subredditl. If there is
// no tip, fetchTip considers the call an adjustment round, and will fetch a new
// reference tip but discard the post (because, most likely, that post was
// already returned before).
func (l *listingScanner) fetchTip() (
	[]*redditproto.Link,
	[]*redditproto.Comment,
	error,
) {
	tip := l.tip[len(l.tip)-1]
	limit := uint(operator.MaxLinks)
	adjustment := false
	if tip == "" {
		limit = 1
		adjustment = true
	}

	links := ([]*redditproto.Link)(nil)
	comments := ([]*redditproto.Comment)(nil)
	err := (error)(nil)

	if l.user != "" {
		links, comments, err = l.op.UserContent(
			l.user,
			"",
			tip,
			limit,
		)
	} else {
		links, err = l.op.Posts(
			l.subreddit,
			"",
			tip,
			limit,
		)
	}
	if err != nil {
		return nil, nil, err
	}

	j := 0
	k := 0
	tipAppendage := make([]string, len(links)+len(comments))
	for i := 0; i < len(links)+len(comments); i++ {
		if k < len(comments) && (j > len(links) || links[j].GetCreatedUtc() > comments[k].GetCreatedUtc()) {
			tipAppendage[i] = comments[k].GetName()
			k++
		} else if j < len(links) {
			tipAppendage[i] = links[j].GetName()
			j++
		}
	}

	l.tip = append(l.tip, tipAppendage...)
	if len(l.tip) > maxTipSize {
		l.tip = l.tip[len(l.tip)-maxTipSize:]
	}

	if adjustment && len(links)+len(comments) == 1 {
		return nil, nil, nil
	}

	return links, comments, nil
}

// fixTip attempts to fix the PostMonitor's reference point for new postl. If it
// has been deleted, fixTip will move to a fallback tip. fixtip returns true if
// the tip was shaved.
func (l *listingScanner) fixTip() (bool, error) {
	exists, err := l.op.IsThereThing(l.tip[len(l.tip)-1])
	if err != nil {
		return false, err
	}

	if exists == false {
		l.shaveTip()
	}

	return !exists, nil
}

// shaveTip shaves off the latest tip thread name. If all tips are shaved off,
// uses an empty tip name (this will just get the latest threads).
func (l *listingScanner) shaveTip() {
	if len(l.tip) == 1 {
		l.tip[0] = ""
		return
	}

	l.tip = l.tip[:len(l.tip)-1]
}
