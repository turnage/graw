package reddit

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
}

func (a App) configured() bool {
	allNotEmpty := func(ss ...string) bool {
		for _, s := range ss {
			if s == "" {
				return false
			}
		}
		return true
	}

	return allNotEmpty(a.tokenURL, a.ID, a.Secret, a.Username, a.Password)
}
