package reddit

import (
	"net/http"
)

// mockClient stores the request it receives.
type mockClient struct {
	request *http.Request
}

func (m *mockClient) Do(r *http.Request) ([]byte, error) {
	m.request = r
	return nil, nil
}
