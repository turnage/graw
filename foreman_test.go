package graw

import (
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/mix/graw/reddit"
)

type mockBot struct {
	err            error
	setUpCalled    bool
	tearDownCalled bool
}

func (m *mockBot) SetUp() error {
	m.setUpCalled = true
	return m.err
}

func (m *mockBot) TearDown() {
	m.tearDownCalled = true
}

func TestForemanControls(t *testing.T) {
	b := &mockBot{}
	errs := make(chan error)
	kill := make(chan bool)
	testLogger := log.New(io.Discard, "", 0)

	stop, _, err := launch(b, kill, errs, testLogger)
	if err != nil {
		t.Fatalf("error launching the foreman: %v", err)
	}

	stop()

	if !b.setUpCalled {
		t.Errorf("SetUp() was not called on bot")
	}

	if !b.tearDownCalled {
		t.Errorf("TearDown() was not called on bot")
	}
}

func TestForemanError(t *testing.T) {
	errs := make(chan error)

	kill := make(chan bool)
	result := make(chan error)
	testLogger := log.New(io.Discard, "", 0)

	_, wait, err := launch(nil, kill, errs, testLogger)
	if err != nil {
		t.Fatalf("error launching the foreman: %v", err)
	}
	go func() {
		result <- wait()
	}()

	// send errors that should be ignored and then a unique error to make
	// sure only the unique error killed the foreman (verified by checking
	// that it is the one that comes through the result channel, since the
	// errors are read chronologically).

	uniqueError := fmt.Errorf("an error")

	go func() {
		errs <- nil
		errs <- reddit.BusyErr
		errs <- reddit.GatewayErr
		errs <- reddit.GatewayTimeoutErr
		errs <- uniqueError
	}()
	waitForForeman(result, uniqueError, t)

}

func waitForForeman(result <-chan error, expected error, t *testing.T) {
	select {
	case err := <-result:
		if err != expected {
			t.Errorf("error from foreman run: %v", err)
		}
	case <-time.After(time.Second):
		t.Errorf("foreman did not stop()")
	}
}
