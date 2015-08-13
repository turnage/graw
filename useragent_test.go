package graw

import (
	"io/ioutil"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/turnage/redditproto"
)

func TestNewUserAgent(t *testing.T) {
	expected := &redditproto.UserAgent{}
	if err := proto.UnmarshalText(`
		user_agent: "test"
		client_id: "id"
		client_secret: "secret"
		username: "user"
		password: "1234"
	`, expected); err != nil {
		t.Errorf("could not build expectation proto: %v", err)
	}

	actual := NewUserAgent("test", "id", "secret", "user", "1234")
	if !proto.Equal(expected, actual) {
		t.Errorf(
			"user agent incorrect; expected %v, got %v",
			expected,
			actual)
	}
}

func TestNewUserAgentFromtFile(t *testing.T) {
	expected := &redditproto.UserAgent{}
	if err := proto.UnmarshalText(`
		user_agent: "test"
		client_id: "id"
		client_secret: "secret"
		username: "user"
		password: "1234"
	`, expected); err != nil {
		t.Errorf("could not build expectation proto: %v", err)
	}

	agentFile, err := ioutil.TempFile("", "user_agent")
	if err != nil {
		t.Errorf("could not make user_agent file: %v", err)
	}

	_, err = agentFile.WriteString(proto.MarshalTextString(expected))
	if err != nil {
		t.Errorf("could not write to user_agent file: %v", err)
	}

	actual, err := NewUserAgentFromFile(agentFile.Name())
	if err != nil {
		t.Errorf("could not build user agent from file: %v", err)
	}

	if !proto.Equal(expected, actual) {
		t.Errorf(
			"user agent incorrect; expected %v, got %v",
			expected,
			actual)
	}
}
