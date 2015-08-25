package client

import (
	"bytes"
	"io"
	"net/http"
)

// rcloser implements ReadCloser for bytes.Buffer.
type rcloser struct {
	*bytes.Buffer
}

func (r *rcloser) Close() error {
	return nil
}

type mockClient struct {
	response []byte
	err      error
}

// Do returns the preconfigured response in the mock client.
func (m *mockClient) Do(r *http.Request) (io.ReadCloser, error) {
	return &rcloser{bytes.NewBuffer(m.response)}, m.err
}
