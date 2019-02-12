package reddit

import (
	"testing"
)

func TestAppUnauthenticated(t *testing.T) {
	for i, test := range []struct {
		input  App
		output bool
	}{
		{App{"y", "", "", "", ""}, true},
		{App{"", "y", "", "", ""}, true},
		{App{"y", "y", "", "", ""}, false},
		{App{"y", "y", "y", "", ""}, false},
		{App{"y", "y", "", "y", ""}, false},
		{App{"y", "y", "y", "y", ""}, false},
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
		{App{"", "", "", "", ""}, errMissingOauthCredentials},
		{App{"y", "", "", "", ""}, errMissingOauthCredentials},
		{App{"", "y", "", "", ""}, errMissingOauthCredentials},
		{App{"y", "y", "y", "", ""}, errMissingPassword},
		{App{"y", "y", "", "y", ""}, errMissingUsername},
		{App{"y", "y", "", "", ""}, nil},
		{App{"y", "y", "y", "y", ""}, nil},
	} {
		if actual := test.input.validateAuth(); actual != test.output {
			t.Errorf("wrong on %d; wanted %v", i, test.output)
		}
	}
}
