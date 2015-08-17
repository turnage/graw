package client

import (
	"io/ioutil"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/turnage/redditproto"
)

func TestNew(t *testing.T) {
	if _, err := New("fakefile"); err == nil {
		t.Errorf("wanted to return error for nonexistent file")
	}

	testInput := &redditproto.UserAgent{}
	if err := proto.UnmarshalText(`
		user_agent: "test"
		client_id: "id"
		client_secret: "secret"
		username: "user"
		password: "pass"
	`, testInput); err != nil {
		t.Errorf("failed to build test expectation proto: %v", err)
	}

	testFile, err := ioutil.TempFile("", "user_agent")
	if err != nil {
		t.Errorf("failed to make test input file: %v", err)
	}

	if err := proto.MarshalText(testFile, testInput); err != nil {
		t.Errorf("failed to write test input file: %v", err)
	}

	if _, err := New(testFile.Name()); err != nil {
		t.Errorf("error: %v", err)
	}
}

func TestNewMock(t *testing.T) {
	resp := `{"key":"value"}`
	expected := &struct {
		Key string
	}{Key: "value"}
	actual := &struct {
		Key string `"json:"key,omitempty"`
	}{}
	mock := NewMock(resp)
	if err := mock.Do(nil, actual); err != nil {
		t.Fatalf("Do() failed: %v", err)
	}
	if actual.Key != expected.Key {
		t.Errorf(
			"response incorrect; got %v, wanted %v",
			actual,
			expected)
	}
}
