package auth

import (
	"net/http"

	"golang.org/x/oauth2"
)

// oauth2Authorizer implements Authorizer for oauth2 logged-in auth.
type oauth2Authorizer struct {
	// id is a unique identifying id for a bot or script.
	id string
	// secret is a unique never-shared string for a bot or script.
	secret string
	// user is a reddit username.
	user string
	// pass is the password for user.
	pass string
}

// NewOAuth2Authorizer returns an Authorizer using logged-in OAuth2.
func NewOAuth2Authorizer(id, secret, user, pass string) Authorizer {
	return &oauth2Authorizer{
		id:     id,
		secret: secret,
		user:   user,
		pass:   pass,
	}
}

// Client returns an *http.Client that automatically refreshes oauth2 tokens
// when they expire, and includes the authorization information in the header of
// all requests.
func (o *oauth2Authorizer) Client(authURL string) (*http.Client, error) {
	conf := &oauth2.Config{
		ClientID:     o.id,
		ClientSecret: o.secret,
		Endpoint: oauth2.Endpoint{
			TokenURL: authURL,
		},
	}

	token, err := conf.PasswordCredentialsToken(
		oauth2.NoContext,
		o.user,
		o.pass)
	if err != nil {
		return nil, err
	}

	return conf.Client(oauth2.NoContext, token), nil
}
