package auth

import (
	"net/http"
)

// mockAuth mocks an implementation of Authorizer.
type mockAuth struct {
	// client will be returned in calls to Client().
	client *http.Client
	// err will be returned in calls to Client().
	err error
}

// NewMockAuth returns a mock implementation of Authorizer that regurgitates the
// values passed here from calls to Client().
func NewMockAuth(client *http.Client, err error) Authorizer {
	return &mockAuth{client: client, err: err}
}

// Client returns the preset return values for the mock implementation.
func (m *mockAuth) Client(authURL string) (*http.Client, error) {
	return m.client, m.err
}
