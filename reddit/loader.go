package reddit

import (
	"bytes"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/turnage/redditproto"
)

// load loads the user agent and App config from an AgentFile (legacy graw 0.3.0
// file format).
func load(filename string) (string, App, error) {
	agentPB, err := loadAgentFile(filename)
	return agentPB.GetUserAgent(), App{
		ID:       agentPB.GetClientId(),
		Secret:   agentPB.GetClientSecret(),
		Username: agentPB.GetUsername(),
		Password: agentPB.GetPassword(),
	}, err
}

// loadAgentFile reads a user agent from a protobuffer file and returns it.
func loadAgentFile(filename string) (*redditproto.UserAgent, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	agent := &redditproto.UserAgent{}
	return agent, proto.UnmarshalText(bytes.NewBuffer(buf).String(), agent)
}
