package auth

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var mock = &mockAuth{
	client: http.DefaultClient,
	err:    fmt.Errorf("BAD STUFF WENT DOWN"),
}

func TestNewMockAuth(t *testing.T) {
	actual := NewMockAuth(mock.client, mock.err).(*mockAuth)
	if !reflect.DeepEqual(actual, mock) {
		t.Errorf(
			"mock built incorrectly; got %v, wanted %v",
			actual,
			mock)
	}
}

func TestClientMockAuth(t *testing.T) {
	client, err := mock.Client("any string")
	if err != mock.err {
		t.Error("mock's preset err not returned")
	}
	if client != mock.client {
		t.Error("mock's preset client not returned")
	}
}
