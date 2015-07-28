package client

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNewMockClient(t *testing.T) {
	expectedError := fmt.Errorf("SOMETHING WENT WRONG")
	expectedResponse := &http.Response{}
	mock := NewMockClient(expectedResponse, expectedError).(*mockClient)

	if mock == nil {
		t.Fatal("no mock client returned")
	}

	if mock.response != expectedResponse {
		t.Errorf(
			"mock response incorrect; got %v, wanted %v",
			mock.response,
			expectedResponse)
	}

	if mock.err != expectedError {
		t.Errorf(
			"mock error incorrect; got %v, wanted %v",
			mock.err,
			expectedError)
	}
}

func TestMockDo(t *testing.T) {
	expectedResp := &http.Response{StatusCode: 200}
	expectedErr := fmt.Errorf("BAD THING")
	mock := &mockClient{response: expectedResp, err: expectedErr}
	actualResp, actualErr := mock.Do(nil)

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
