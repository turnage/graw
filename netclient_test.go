package graw

import (
	"net/http"
	"testing"
)

func TestNetDo(t *testing.T) {
	client := &netClient{Client: http.DefaultClient}
	server := newServerFromResponse(200, []byte("Hello"))
	request, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Error("error building request")
	}
	resp, err := client.Do(request)
	if err != nil {
		t.Error("error executing request")
	}
	if !responseIs(resp, 200, []byte("Hello")) {
		t.Error("error in request")
	}
}
