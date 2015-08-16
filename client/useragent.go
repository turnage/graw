package client

import (
	"bytes"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/turnage/redditproto"
)

// load reads a user agent from a protobuffer file and returns it.
func load(filename string) (*redditproto.UserAgent, error) {
	agentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(agentBytes)
	agent := &redditproto.UserAgent{}
	if err := proto.UnmarshalText(buffer.String(), agent); err != nil {
		return nil, err
	}

	return agent, nil
}
