// Package monitor includes monitors for different parts of Reddit, such as a
// user inbox or a subreddit's post feed.
package monitor

// Monitor defines the controls for a Monitor.
type Monitor interface {
	// Update will check for new events, and send them to the Monitor's
	// handler.
	Update() error
}
