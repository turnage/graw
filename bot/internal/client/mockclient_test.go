package client

import (
	"fmt"
	"net/http"
	"testing"
)

func TestMockDo(t *testing.T) {
	err := fmt.Errorf("a real bad thing")
	resp := &http.Response{Status: "pretty ok, how about you?"}
	mock := &mockClient{
		resp: resp,
		err:  err,
	}
	actualResp, actualErr := mock.Do(nil)
	if actualResp != resp {
		t.Errorf("wanted %v; got %v", resp, actualResp)
	}
	if actualErr != err {
		t.Errorf("wanted %v; got %v", err, actualErr)
	}

}
