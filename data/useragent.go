package data

import (
	"bytes"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

// NewUserAgent returns a new UserAgent containing the provided fields.
func NewUserAgent(userAgent, id, secret, user, pass string) *UserAgent {
	return &UserAgent{
		UserAgent:    &userAgent,
		ClientId:     &id,
		ClientSecret: &secret,
		Username:     &user,
		Password:     &pass,
	}
}

// NewUserAgent returns a new UserAgent from a protobuffer file.
func NewUserAgentFromFile(filename string) (*UserAgent, error) {
	agentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	agentText := bytes.NewBuffer(agentBytes)
	agent := &UserAgent{}
	if err := proto.UnmarshalText(agentText.String(), agent); err != nil {
		return nil, err
	}

	return agent, nil
}

