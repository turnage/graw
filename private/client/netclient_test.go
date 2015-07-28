package client

import (
	"net/http"
	"testing"

	"github.com/paytonturnage/graw/private/testutil"
)

func TestNewNetClient(t *testing.T) {
	client := NewNetClient(http.DefaultClient).(*netClient)
	if client.client != http.DefaultClient {
		t.Errorf(
			"client incorrect; got %v, wanted %v",
			client.client,
			http.DefaultClient)
	}
}

func TestNetDo(t *testing.T) {
	client := &netClient{client: http.DefaultClient}
	server := testutil.NewServerFromResponse(200, []byte("Hello"))
	request, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Error("error building request")
	}
	resp, err := client.Do(request)
	if err != nil {
		t.Error("error executing request")
	}
	if !testutil.ResponseIs(resp, 200, []byte("Hello")) {
		t.Error("error in request")
	}
}
