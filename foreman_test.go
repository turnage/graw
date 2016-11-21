package graw

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/turnage/graw/reddit"
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
	result, stop := testForeman(b, nil, t)

	stop()
	waitForForeman(result, nil, t)

	if !b.setUpCalled {
		t.Errorf("SetUp() was not called on bot")
	}

	if !b.tearDownCalled {
		t.Errorf("TearDown() was not called on bot")
	}
}

func TestForemanError(t *testing.T) {
	errs := make(chan error)
	result, _ := testForeman(nil, errs, t)

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

func testForeman(handler interface{}, errs chan error, t *testing.T) (
	<-chan error,
	func(),
) {
	kill := make(chan bool)
	result := make(chan error)
	logger := log.New(ioutil.Discard, "", 0)
	if errs == nil {
		errs = make(chan error)
	}

	stop, wait, err := launch(handler, kill, errs, logger)
	if err != nil {
		t.Fatalf("error launching the foreman: %v", err)
	}

	go func() {
		result <- wait()
	}()

	return result, stop
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
