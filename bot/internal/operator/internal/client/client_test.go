package client

import (
	"bytes"
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
	expected := "internet"
	mock := NewMock(expected)
	body, err := mock.Do(nil)
	if err != nil {
		t.Fatalf("Do() error: %v", err)
	}

	actualBuffer := new(bytes.Buffer)
	_, err = actualBuffer.ReadFrom(body)
	if err != nil {
		t.Errorf("failed to read response body: %v", err)
	}

	actual := actualBuffer.String()
	if actual != expected {
		t.Errorf("got %v, wanted %v", actual, expected)
	}
}
