package client

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

func TestNewAppClient(t *testing.T) {
	serv := serverWhich([]byte(`{
		"access_token": "aksjdsakjd",
		"token_type": "bearer",
		"expires_in": 100,
		"scope": "*",
		"refresh_token": "sidfnsidfnsd" 
	}`), http.StatusOK)

	if client, err := New(
		Config{
			App: App{
				TokenURL: serv.URL,
				ID:       "id",
				Secret:   "secret",
				Username: "user",
				Password: "password",
			},
		},
	); err != nil {
		t.Errorf("failed to fetch token: %v", err)
	} else if client == nil {
		t.Errorf("client was nil")
	} else if app, ok := client.(*appClient); !ok {
		t.Errorf("client was not an appClient")
	} else if app.token == nil {
		t.Errorf("appClient's token was not set")
	}
}

func TestNewAnonClient(t *testing.T) {
	if client, err := New(Config{}); err != nil {
		t.Errorf("error making anon client")
	} else if client == nil {
		t.Errorf("anon client was nil")
	} else if _, ok := client.(*base); !ok {
		t.Errorf("anon client was not a base implementation")
	}
}

func TestDo(t *testing.T) {
	r := &base{cli: &http.Client{}}
	for _, test := range []struct {
		body []byte
		code int
		err  error
	}{
		{[]byte("expected"), http.StatusOK, nil},
		{nil, http.StatusForbidden, PermissionDeniedErr},
		{nil, http.StatusServiceUnavailable, BusyErr},
		{nil, http.StatusTooManyRequests, RateLimitErr},
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
