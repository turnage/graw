package client

import (
	"bytes"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/turnage/redditproto"
)

// load reads a user agent from a protobuffer file and returns it.
func load(filename string) (*redditproto.UserAgent, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	agent := &redditproto.UserAgent{}
	return agent, proto.UnmarshalText(bytes.NewBuffer(buf).String(), agent)
}
