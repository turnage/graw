package graw

import (
	"fmt"
	"net/http"
	"testing"
)

func TestMockDo(t *testing.T) {
	expectedResp := &http.Response{StatusCode: 200}
	expectedErr := fmt.Errorf("BAD THING")
	mock := &mockClient{response: expectedResp, err: expectedErr}
	actualResp, actualErr := mock.do(nil)

	if actualErr != expectedErr {
		t.Errorf(
			"err incorrect; got %v, wanted %v",
			actualErr,
			expectedErr)
	}

	if actualResp != expectedResp {
		t.Errorf(
			"resp incorrect; got %v, wanted %v",
			actualResp,
			expectedResp)
	}
}
