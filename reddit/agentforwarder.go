package reddit

import (
	"net/http"
)

// agentForward forwards a user agent in all requests made by the Transport.
type agentForwarder struct {
	http.Transport
	agent string
}

// RoundTrip sets a predefined agent in the request and then forwards it to the
// default RountTrip implementation.
func (a *agentForwarder) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("User-Agent", a.agent)
	return a.Transport.RoundTrip(r)
}

func clientWithAgent(agent string) *http.Client {
	return &http.Client{Transport: &agentForwarder{agent: agent}}
}
