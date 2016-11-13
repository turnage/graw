package engine

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/turnage/graw/internal/dispatcher"
)

type mockDispatcher struct {
	err   error
	calls chan bool
}

func (m *mockDispatcher) Dispatch() error {
	m.calls <- true
	return m.err
}

func TestRun(t *testing.T) {
	d1 := &mockDispatcher{nil, make(chan bool)}
	d2 := &mockDispatcher{nil, make(chan bool)}
	rate := make(chan time.Time)
	logger := log.New(os.Stderr, "", log.LstdFlags)
	e := New(Config{[]dispatcher.Dispatcher{d1, d2}, rate, logger})

	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = e.Run()
		wg.Done()
	}()

	// Allow it to call the first dispatcher.
	rate <- time.Now()
	<-d1.calls

	// Allow it to call the second dispatcher.
	rate <- time.Now()
	<-d2.calls

	// It should call the first dispatcher again.
	rate <- time.Now()
	<-d1.calls

	e.Stop()
	wg.Wait()

	if err != nil {
		t.Errorf("error running engine: %v", err)
	}
}

func TestErrBubble(t *testing.T) {
	expectedErr := fmt.Errorf("an error")
	d := &mockDispatcher{expectedErr, make(chan bool)}
	rate := make(chan time.Time)
	logger := log.New(os.Stderr, "", log.LstdFlags)
	e := New(Config{[]dispatcher.Dispatcher{d}, rate, logger})

	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = e.Run()
		wg.Done()
	}()

	rate <- time.Now()
	<-d.calls

	wg.Wait()

	if err != expectedErr {
		t.Errorf("got %v; wanted %v", err, expectedErr)
	}
}
