package reddit

import (
	"io/ioutil"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/turnage/redditproto"
)

func TestLoad(t *testing.T) {
	expected := &redditproto.UserAgent{}
	if err := proto.UnmarshalText(`
		user_agent: "test"
		client_id: "id"
		client_secret: "secret"
		username: "user"
		password: "pass"
	`, expected); err != nil {
		t.Errorf("Failed to build test expectation proto: %v", err)
	}

	testFile, err := ioutil.TempFile("", "user_agent")
	if err != nil {
		t.Errorf("Failed to make test input file: %v", err)
	}

	if err := proto.MarshalText(testFile, expected); err != nil {
		t.Errorf("Failed to write test input file: %v", err)
	}

	if _, err := loadAgentFile("notarealfile"); err == nil {
		t.Error("Wanted error returned with nonexistent file as input")
	}

	actual, err := loadAgentFile(testFile.Name())
	if err != nil {
		t.Errorf("Failed: %v", err)
	}

	if !proto.Equal(expected, actual) {
		t.Errorf("Got %v; wanted %v", actual, expected)
	}
}
