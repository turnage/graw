package reddit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientWithAgent(t *testing.T) {
	c := clientWithAgent("agent")

	forwarder := c.Transport
	if v, ok := forwarder.(*agentForwarder); ok {
		if v.agent != "agent" {
			t.Errorf("expected `agent`, got %s", v.agent)
		}
	} else {
		t.Error("expected *agentForwarder")
	}
}

type mockTransport struct{}

func (m mockTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, nil
}

func TestPatchWithAgent(t *testing.T) {
	t.Run("null-transport", func(t *testing.T) {
		client := &http.Client{}
		c := patchWithAgent(client, "agent")

		forwarder := c.Transport
		if v, ok := forwarder.(*agentForwarder); ok {
			if v.agent != "agent" {
				t.Errorf("expected `agent`, got %s", v.agent)
			}
		} else {
			t.Error("expected *agentForwarder")
		}
	})

	t.Run("predefined-transport", func(t *testing.T) {
		client := &http.Client{
			Transport: mockTransport{},
		}
		c := patchWithAgent(client, "agent")

		forwarder := c.Transport
		if _, ok := forwarder.(mockTransport); ok {
			t.Error("expected mockTransport")
		}
	})
}

func TestAgentForwarder_RoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if v := req.Header.Get("User-Agent"); v != "agent" {
			t.Errorf("expected `agent`, got %s", v)
		}
	}))

	defer server.Close()

	client := patchWithAgent(server.Client(), "agent")

	_, err := client.Get("https://example.com")
	if err != nil {
		t.Error(err)
		return
	}
}
