package graw

import (
	"testing"
)

func TestClientOAuth2(t *testing.T) {
	serv := newServerFromResponse(200, []byte(`{
			"access_token": "sjkhefwhf383nfjkf",
			"token_type": "bearer",
			"expires_in": 3600
			"scope": "*",
			"refresh_token": "akjfbkfjhksdjhf"
	}`))
	client, err := oauth("", "", "", "", serv.URL)
	if err != nil {
		t.Errorf("failed to make oauth client: %v", err)
	}

	if client == nil {
		t.Error("client not returned")
	}

	serv = newServerFromResponse(404, []byte("404"))
	if _, err := oauth("", "", "", "", serv.URL); err == nil {
		t.Error("bad response did not generate error")
	}
}
