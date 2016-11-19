package streams

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/turnage/graw/reddit"
)

type mockMonitor struct {
	h   reddit.Harvest
	err error
}

func (m *mockMonitor) Update() (reddit.Harvest, error) {
	return m.h, m.err
}

func TestStream(t *testing.T) {
	kill := make(chan bool)
	errs := make(chan error)
	mon := &mockMonitor{
		h: reddit.Harvest{
			Posts: []*reddit.Post{&reddit.Post{Title: "Title"}},
			Comments: []*reddit.Comment{
				&reddit.Comment{Body: "body"},
			},
			Messages: []*reddit.Message{
				&reddit.Message{Body: "body"},
			},
		},
	}

	posts, comments, messages := stream(mon, kill, errs)

	done := make(chan bool)
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		<-posts
		wg.Done()
	}()
	go func() {
		<-comments
		wg.Done()
	}()
	go func() {
		<-messages
		wg.Done()
	}()
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		close(errs)
	case <-time.After(time.Second):
		t.Errorf("stream did not emit all the items expected from it")
	}
}

func TestErrorPropagation(t *testing.T) {
	done := make(chan bool)
	errs := make(chan error)
	kill := make(chan bool)
	posts := make(chan *reddit.Post)
	comments := make(chan *reddit.Comment)
	messages := make(chan *reddit.Message)
	mon := &mockMonitor{err: fmt.Errorf("an error")}
	go func() {
		flow(mon, kill, errs, posts, comments, messages)
		done <- true
	}()
	go func() {
		<-errs
		// Spawn a routine to consume errors once we know flow is
		// emitting them, so it doesn't get blocked.
		go func() {
			for range errs {
			}
		}()
		close(kill)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Errorf("loop did not report error or accept kill")
	}
}
