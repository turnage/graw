package reaper

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func serverWhich(body string, code int) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
				w.Write([]byte(body))
			},
		),
	)
}

func TestNew(t *testing.T) {
	serv := serverWhich(`{
		"access_token": aksjdsakjd,
		"token_type": "bearer",
		"expires_in": 100,
		"scope": "*",
		"refresh_token": sidfnsidfnsd 
	}`, http.StatusOK)

	if reaper, err := New(Config{TokenURL: serv.URL}); err != nil {
		t.Errorf("failed to fetch token: %v", err)
	} else if reaper == nil {
		t.Errorf("reaper was nil")
	}
}

func TestReap(t *testing.T) {
	r := &reaper{cli: &http.Client{}}
	for _, test := range []struct {
		body string
		code int
		err  error
	}{
		{"expected", http.StatusOK, nil},
		{"", http.StatusForbidden, PermissionDeniedErr},
		{"", http.StatusServiceUnavailable, BusyErr},
		{"", http.StatusTooManyRequests, RateLimitErr},
		{"", http.StatusOK, nil},
	} {
		serv := serverWhich(test.body, test.code)

		req, err := http.NewRequest("GET", serv.URL, nil)
		if err != nil {
			t.Fatalf("failed to prepare request for test: %v", err)
		}

		if body, err := r.reap(req); err != test.err {
			t.Errorf("unexpected error: %v", err)
		} else if reflect.DeepEqual(body, test.body) == false {
			t.Errorf("body got %s; wanted %s", body, test.body)
		}
	}
}
