package reddit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func serverWhich(body []byte, code int) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
				w.Write(body)
			},
		),
	)
}

func TestNewAnonClient(t *testing.T) {
	if client, err := newClient(clientConfig{}); err != nil {
		t.Error("error making anon client")
	} else if client == nil {
		t.Error("anon client was nil")
	} else if _, ok := client.(*baseClient); !ok {
		t.Error("anon client was not a base implementation")
	}
}

func TestDo(t *testing.T) {
	r := &baseClient{cli: &http.Client{}}
	for _, test := range []struct {
		body []byte
		code int
		err  error
	}{
		{[]byte("expected"), http.StatusOK, nil},
		{nil, http.StatusForbidden, ErrPermissionDenied},
		{nil, http.StatusServiceUnavailable, ErrBusy},
		{nil, http.StatusTooManyRequests, ErrRateLimit},
		{nil, http.StatusBadGateway, ErrBadGateway},
		{nil, http.StatusGatewayTimeout, ErrGatewayTimeout},
		{nil, http.StatusOK, nil},
	} {
		serv := serverWhich(test.body, test.code)

		req, err := http.NewRequest("GET", serv.URL, nil)
		if err != nil {
			t.Fatalf("failed to prepare request for test: %v", err)
		}

		body, err := r.Do(req)
		if err != test.err {
			t.Errorf("unexpected error: %v", err)
		} else if len(body) != len(test.body) {
			t.Errorf(
				"unexpected body length; got %d and wanted %d",
				len(body), len(test.body),
			)
		}

		for i := 0; i < len(body); i++ {
			if body[i] != test.body[i] {
				t.Errorf(
					"body got %s; wanted %s",
					body, test.body,
				)
				break
			}
		}
	}
}
