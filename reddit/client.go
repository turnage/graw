package reddit

import (
	"net/http"
)

// tokenURL is the url of reddit's oauth2 authorization service.
const tokenURL = "https://www.reddit.com/api/v1/access_token"

// clientConfig holds all the information needed to define Client behavior, such
// as who the client will identify as externally and where to authorize.
type clientConfig struct {
	// Agent is the user agent set in all requests made by the Client.
	agent string

	// If all fields in App are set, this client will attempt to identify as
	// a registered Reddit app using the credentials.
	app App
}

// client executes http Requests and invisibly handles OAuth2 authorization.
type client interface {
	Do(*http.Request) ([]byte, error)
	AutoRefresh() error
}

// newClient returns a new client using the given user to make requests.
func newClient(c clientConfig) (client, error) {
	if c.app.tokenURL == "" {
		c.app.tokenURL = tokenURL
	}

	if c.app.unauthenticated() {
		return &appClient{
			app: &c.app,
		}, nil
	}

	if err := c.app.validateAuth(); err != nil {
		return nil, err
	}

	return newAppClient(c)
}
