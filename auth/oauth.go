package auth

import (
	"errors"
	"net/url"
	"time"

	"github.com/paytonturnage/graw/nface"
)

// mode describes OAuth modes.
type mode int

const (
	// APP is the app-only OAuth mode.
	APP = iota
	// USER is the user-based OAuth mode.
	USER = iota
)

const (
	// appGrantType is the app-only OAuth grant type.
	appGrantType = "https://oauth.reddit.com/grants/installed_client"
	// baseUrl is the api url for requesting OAuth tokens.
	baseUrl = "https://www.reddit.com/api/v1/access_token"
	// timeBuffer is the time before actual token expiry to refresh (secs).
	timeBuffer = 30
	// userGrantType is the user-based OAuth grant type.
	userGrantType = "password"
	// userRefreshType is the user-based OAuth refresh token grant type.
	userRefreshType = "refresh_token"
)

// OAuth describes and manages an OAuth identity.
type OAuth struct {
	// appID is the application id (see comments on NewOAuth).
	appID string
	// appSecret is the application secret (see comments on NewOAuth).
	appSecret string
	// acctUser is the account username for logged-in OAuths.
	acctUser string
	// acctPass is the account password for logged-in OAuths.
	acctPass string
	// device is the device id (uuid) of the session (app-only oauth).
	device string
	// expiry is the time at which this token expires.
	expiry time.Time
	// oauthMode indicates the mode of auth (app-only or through a user).
	oauthMode mode
	// scope is the scope for which this token grants access. Search this
	// page for "scope" for more information:
	// https://www.reddit.com/dev/api/oauth
	scope string
	// token is the access token to authenticate via OAuth.
	token string
}

// oauthGrant describes the JSON OAuth response from the reddit API.
type oauthGrant struct {
	// AccessToken is the oauth identity.
	AccessToken string `json:"access_token"`
	// Error contains an error message if the grant failed.
	Error string `json:"error"`
	// ExpiresIn is the seconds from grant that the token is valid for.
	ExpiresIn int `json:"expires_in"`
	// Scope is the permissions scope of this token.
	Scope string `json:"scope"`
}

// NewAppOAuth returns an OAuth struct which will authenticate app-only (with no
// association to a user account).
func NewAppOAuth(id, secret, device string) *OAuth {
	return &OAuth{
		appID:     id,
		appSecret: secret,
		device:    device,
		oauthMode: APP,
	}
}

// NewUserOAuth returns an OAuth struct which will authenticate through a user
// account.
func NewUserOAuth(id, secret, user, pass string) *OAuth {
	return &OAuth{
		appID:     id,
		appSecret: secret,
		acctUser:  user,
		acctPass:  pass,
		oauthMode: USER,
	}
}

// Token returns an OAuth token. If no error is returned, this token will be
// valid and up to date.
func (o *OAuth) Token() (string, error) {
	if time.Now().Before(o.expiry) {
		return o.token, nil
	}

	if o.oauthMode == APP {
		if err := o.newApp(); err != nil {
			return "", err
		}
		return o.Token()
	}

	if err := o.refreshUser(); err != nil {
		return "", err
	}
	return o.Token()
}

// updateGrant updates the OAuth's grant information with a new grant.
func (o *OAuth) updateGrant(resp *oauthGrant) {
	o.token = resp.AccessToken
	o.expiry = time.Now().Add(
		time.Duration(resp.ExpiresIn-timeBuffer) * time.Second)
	o.scope = resp.Scope
}

// newApp requests a new app-only OAuth token.
func (o *OAuth) newApp() error {
	resp, err := oauthRequest(&nface.Request{
		Action:        nface.POST,
		BasicAuthUser: o.appID,
		BasicAuthPass: o.appSecret,
		BaseUrl:       baseUrl,
		Values: &url.Values{
			"grant_type": []string{appGrantType},
			"device_id":  []string{o.device},
		}})
	if err != nil {
		return err
	}

	o.updateGrant(resp)

	return nil
}

// newUser requests a new user-based OAuth token.
func (o *OAuth) newUser() error {
	resp, err := oauthRequest(&nface.Request{
		Action:        nface.POST,
		BasicAuthUser: o.appID,
		BasicAuthPass: o.appSecret,
		BaseUrl:       baseUrl,
		Values: &url.Values{
			"grant_type": []string{userGrantType},
			"username":   []string{o.acctUser},
			"password":   []string{o.acctPass},
		}})
	if err != nil {
		return err
	}

	o.updateGrant(resp)

	return nil
}

// refreshUser refreshes a user-based OAuth token. If there is no token to
// refresh, requests a new one.
func (o *OAuth) refreshUser() error {
	if o.token == "" {
		return o.newUser()
	}

	resp, err := oauthRequest(&nface.Request{
		Action:        nface.POST,
		BasicAuthUser: o.appID,
		BasicAuthPass: o.appSecret,
		BaseUrl:       baseUrl,
		Values: &url.Values{
			"grant_type":    []string{userRefreshType},
			"username":      []string{o.acctUser},
			"password":      []string{o.acctPass},
			"duration":      []string{"permanent"},
			"refresh_token": []string{o.token},
		}})
	if err != nil {
		return err
	}

	o.updateGrant(resp)

	return nil
}

// oauthRequest executes a Request and parses the response into an oauthGrant.
func oauthRequest(req *nface.Request) (*oauthGrant, error) {
	resp := &oauthGrant{}
	if err := nface.Exec(req, resp); err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	return resp, nil
}
