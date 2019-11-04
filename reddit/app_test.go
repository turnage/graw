package reddit

import (
	"testing"

	"golang.org/x/oauth2"
)

func TestAppUnauthenticated(t *testing.T) {
	for i, test := range []struct {
		input  App
		output bool
	}{
		{App{"y", "", "", "", "", nil}, true},
		{App{"", "y", "", "", "", nil}, true},
		{App{"y", "", "", "", "", &oauth2.Token{}}, false},
		{App{"", "y", "", "", "", &oauth2.Token{}}, false},
		{App{"y", "y", "", "", "", nil}, false},
		{App{"y", "y", "y", "", "", nil}, false},
		{App{"y", "y", "", "y", "", nil}, false},
		{App{"y", "y", "y", "y", "", nil}, false},
	} {
		if actual := test.input.unauthenticated(); actual != test.output {
			t.Errorf("wrong on %d; wanted %v", i, test.output)
		}
	}
}

func TestAppValidateAuth(t *testing.T) {
	for i, test := range []struct {
		input  App
		output error
	}{
		{App{"", "", "", "", "", nil}, errMissingOauthCredentials},
		{App{"y", "", "", "", "", nil}, errMissingOauthCredentials},
		{App{"", "y", "", "", "", nil}, errMissingOauthCredentials},
		{App{"y", "y", "y", "", "", nil}, errMissingPassword},
		{App{"y", "y", "", "y", "", nil}, errMissingUsername},
		{App{"", "", "", "", "", &oauth2.Token{}}, nil},
		{App{"y", "", "", "", "", &oauth2.Token{}}, nil},
		{App{"", "y", "", "", "", &oauth2.Token{}}, nil},
		{App{"y", "y", "y", "", "", &oauth2.Token{}}, nil},
		{App{"y", "y", "", "y", "", &oauth2.Token{}}, nil},
		{App{"y", "y", "", "", "", nil}, nil},
		{App{"y", "y", "y", "y", "", nil}, nil},
	} {
		if actual := test.input.validateAuth(); actual != test.output {
			t.Errorf("wrong on %d; wanted %v", i, test.output)
		}
	}
}
