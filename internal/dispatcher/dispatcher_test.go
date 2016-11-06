package dispatcher

import (
	"fmt"
	"testing"

	"github.com/turnage/graw/internal/data"
	"github.com/turnage/graw/internal/handlers"
	"github.com/turnage/graw/internal/reap"
)

type mockMonitor struct {
	harvest reap.Harvest
}

func (m *mockMonitor) Update() (reap.Harvest, error) {
	return m.harvest, nil
}

func TestDispatch(t *testing.T) {
	receivedPost := false
	receivedComment := false
	receivedMessage := false

	d := New(
		Config{
			Monitor: &mockMonitor{
				reap.Harvest{
					Posts: []*data.Post{
						&data.Post{},
					},
					Comments: []*data.Comment{
						&data.Comment{},
					},
					Messages: []*data.Message{
						&data.Message{},
					},
				},
			},
			PostHandler: handlers.PostHandlerFunc(
				func(_ *data.Post) error {
					receivedPost = true
					return nil
				},
			),
			CommentHandler: handlers.CommentHandlerFunc(
				func(_ *data.Comment) error {
					receivedComment = true
					return nil
				},
			),
			MessageHandler: handlers.MessageHandlerFunc(
				func(_ *data.Message) error {
					receivedMessage = true
					return nil
				},
			),
		},
	)

	if err := d.Dispatch(); err != nil {
		t.Errorf("Error dispatching: %v", err)
	}

	if !receivedPost {
		t.Errorf("Did not received post from dispatcher!")
	}

	if !receivedComment {
		t.Errorf("Did not received comment from dispatcher!")
	}

	if !receivedMessage {
		t.Errorf("Did not received message from dispatcher!")
	}
}

func TestErrorBubble(t *testing.T) {
	expectedErr := fmt.Errorf("expected error")
	d := New(
		Config{
			Monitor: &mockMonitor{
				reap.Harvest{Posts: []*data.Post{&data.Post{}}},
			},
			PostHandler: handlers.PostHandlerFunc(
				func(_ *data.Post) error {
					return expectedErr
				},
			),
		},
	)

	if err := d.Dispatch(); err != expectedErr {
		t.Errorf("Error bubble got %v; wanted %v", err, expectedErr)
	}
}
