package graw

// redditorResponse holds a Redditor struct and any potential error the api call
// expected to return a redditor might return instead.
type redditorResponse struct {
	Redditor
	err string `json:"error,omitempty"`
}
