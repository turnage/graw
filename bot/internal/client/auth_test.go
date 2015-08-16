package client

import (
	"testing"
)

func TestBuild(t *testing.T) {
	if cli := build("id", "secret", "refresh"); cli == nil {
		t.Errorf("did not return client")
	}
}
