package auth

import (
	"reflect"
	"testing"

	"github.com/paytonturnage/graw/private/user/internal/testutil"
)

func TestNewOAuth2Authorizer(t *testing.T) {
	expected := &oauth2Authorizer{
		id:     "id",
		secret: "secret",
		user:   "user",
		pass:   "pass",
	}
	actual := NewOAuth2Authorizer(
		"id",
		"secret",
		"user",
		"pass").(*oauth2Authorizer)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(
			"authorizer made wrong; got %v, wanted %v",
			actual,
			expected)
	}
}

func TestClientOAuth2(t *testing.T) {
	serv := testutil.NewServerFromResponse(200, []byte(`{
			"access_token": "sjkhefwhf383nfjkf",
			"token_type": "bearer",
			"expires_in": 3600
			"scope": "*",
			"refresh_token": "akjfbkfjhksdjhf"
	}`))
	authorizer := &oauth2Authorizer{
		id:     "id",
		secret: "secret",
		user:   "user",
		pass:   "pass",
	}

	client, err := authorizer.Client(serv.URL)
	if err != nil {
		t.Errorf("failed to make oauth client: %v", err)
	}

	if client == nil {
		t.Error("client not returned")
	}

	serv = testutil.NewServerFromResponse(404, []byte("404"))
	if _, err := authorizer.Client(serv.URL); err == nil {
		t.Error("bad response did not generate error")
	}
}
