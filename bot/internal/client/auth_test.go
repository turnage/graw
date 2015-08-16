package client

import (
	"testing"
)

func TestBuild(t *testing.T) {
	if client := build("id", "secret", "refresh"); client == nil {
		t.Errorf("did not return client")
	}
}
