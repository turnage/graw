package reddit

import (
	"testing"
)

func TestAppConfigured(t *testing.T) {
	for i, test := range []struct {
		input  App
		output bool
	}{
		{App{"", "", "", "", ""}, false},
		{App{"y", "y", "y", "y", "y"}, true},
		{App{"", "", "y", "y", "y"}, false},
		{App{"", "", "y", "", ""}, false},
		{App{"", "y", "", "y", ""}, false},
	} {
		if actual := test.input.configured(); actual != test.output {
			t.Errorf("wrong on %d; wanted %v", i, test.output)
		}
	}
}
