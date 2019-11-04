package reddit

import (
	"fmt"

	"golang.org/x/oauth2"
)

var (
	errMissingOauthCredentials = fmt.Errorf("missing oauth credentials")
	errMissingUsername         = fmt.Errorf("missing username")
	errMissingPassword         = fmt.Errorf("missing password")
)

// App holds all the information needed to identify as a registered app on
// Reddit. If you are unfamiliar with this information, you can find it in your
// "apps" tab on reddit; see this tutorial:
// https://github.com/reddit/reddit/wiki/OAuth2
type App struct {
	// ID and Secret are used to claim an OAuth2 grant the bot's account
	// previously authorized.
	ID     string
	Secret string

	// Username and Password are used to authorize with the endpoint.
	Username string
	Password string

	// tokenURL is the url of the token request location for OAuth2.
	tokenURL string

	// If token is specified, username/password authentication is skipped
	Token *oauth2.Token
}

func (a App) unauthenticated() bool {
	return a.Token == nil && (a.ID == "" || a.Secret == "")
}

func (a App) validateAuth() error {
	if a.unauthenticated() {
		return errMissingOauthCredentials
	}

	if a.Token != nil {
		return nil
	}

	if a.Password != "" && a.Username == "" {
		return errMissingUsername
	}

	if a.Username != "" && a.Password == "" {
		return errMissingPassword
	}

	return nil
}
