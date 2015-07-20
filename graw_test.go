package graw

import (
	"encoding/json"
	"testing"

	"github.com/paytonturnage/graw/data"
	"github.com/paytonturnage/graw/nface"
	"github.com/paytonturnage/graw/testutil"
)

func TestMe(t *testing.T) {
	resp, err := json.Marshal(&data.Account{})
	if err != nil {
		t.Fatalf("preparing response failed: %v", err)
	}
	server, serverURL := testutil.NewServerFromResponse(200, resp)
	defer server.Close()

	agent := &Graw{
		client: nface.TestClient(testutil.NewProxyClient(serverURL)),
	}
	if _, err := agent.Me(); err != nil {
		t.Fatalf("failed to get self: %v", err)
	}
}

func TestMeKarma(t *testing.T) {
	resp, err := json.Marshal(&data.KarmaList{})
	if err != nil {
		t.Fatalf("preparing response failed: %v", err)
	}
	server, serverURL := testutil.NewServerFromResponse(200, resp)
	defer server.Close()

	agent := &Graw{
		client: nface.TestClient(testutil.NewProxyClient(serverURL)),
	}
	if _, err := agent.MeKarma(); err != nil {
		t.Fatalf("failed to get karma: %v", err)
	}
}

func TestMeUser(t *testing.T) {
	resp, err := json.Marshal(&data.Account{})
	if err != nil {
		t.Fatalf("preparing response failed: %v", err)
	}
	server, serverURL := testutil.NewServerFromResponse(200, resp)
	defer server.Close()

	agent := &Graw{
		client: nface.TestClient(testutil.NewProxyClient(serverURL)),
	}
	if _, err := agent.User("user"); err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
}
