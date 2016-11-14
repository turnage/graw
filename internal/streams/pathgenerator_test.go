package streams

import (
	"reflect"
	"testing"
)

func TestSubredditsPath(t *testing.T) {
	for i, test := range []struct {
		in  []string
		out string
	}{
		{[]string{"1", "2"}, "/r/1+2/new"},
		{[]string{"1"}, "/r/1/new"},
	} {
		if out := subredditsPath(test.in); out != test.out {
			t.Errorf("%d: got %s; wanted %s", i, out, test.out)
		}
	}
}

func TestUserPaths(t *testing.T) {
	for i, test := range []struct {
		in  []string
		out []string
	}{
		{[]string{"1", "2"}, []string{"/u/1", "/u/2"}},
		{[]string{}, []string{}},
	} {
		if out := userPaths(
			test.in,
		); !reflect.DeepEqual(out, test.out) {
			t.Errorf("%d: got %s; wanted %s", i, out, test.out)
		}
	}
}

func TestLogPathOut(t *testing.T) {
	for i, test := range []struct {
		in  string
		out string
	}{
		{"1", "1.json"},
	} {
		if out := logPathsOut(
			test.in,
		); !reflect.DeepEqual(out, test.out) {
			t.Errorf("%d: got %s; wanted %s", i, out, test.out)
		}
	}
}
