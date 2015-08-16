package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

	if _, err := New(testFile.Name()); err == nil {
		t.Errorf("wanted error for missing refresh token")
	}

	refreshToken := "sldhfslkdjf"
	testInput.RefreshToken = &refreshToken
	testFile, err = ioutil.TempFile("", "user_agent")
	if err != nil {
		t.Errorf("failed to make test input file: %v", err)
	}
	if err := proto.MarshalText(testFile, testInput); err != nil {
		t.Errorf("failed to update test input file: %v", err)
	}

	if _, err := New(testFile.Name()); err != nil {
		t.Errorf("failed to build client: %v", err)
	}
}

func TestNewMock(t *testing.T) {
	err := fmt.Errorf("a real bad thing")
	resp := &http.Response{Status: "pretty ok, how about you?"}
	mock := NewMock(resp, err)
	actualResp, actualErr := mock.Do(nil)
	if actualResp != resp {
		t.Errorf("wanted %v; got %v", resp, actualResp)
	}
	if actualErr != err {
		t.Errorf("wanted %v; got %v", err, actualErr)
	}
}
