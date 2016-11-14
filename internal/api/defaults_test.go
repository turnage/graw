package api

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestWithDefaults(t *testing.T) {
	if diff := pretty.Compare(
		WithDefaults(nil),
		defaultValues,
	); diff != "" {
		t.Errorf("output for nil input wrong; diff: %s", diff)
	}

	if diff := pretty.Compare(
		WithDefaults(map[string]string{"key": "value"}),
		map[string]string{
			"key":      "value",
			"raw_json": "1",
		},
	); diff != "" {
		t.Errorf("output for nonnil input wrong; diff: %s", diff)
	}
}
