package data

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestAccount(t *testing.T) {
	expected := &Account{}
	err := proto.UnmarshalText(`
		name: "fooBar"
		id: "5sryd"
		link_karma: 31
		comment_karma: 557
		is_mod: true
		has_mod_mail: true
		has_mail: true
		inbox_count: 2
		created: 123456789.0
		created_utc: 1315269998.0
		is_gold: true
		gold_credits: 45
		gold_expiration: 23273468.0
		over_18: true
		has_verified_email: true
		hide_from_robots: true
	`, expected)
	actual := &Account{}
	err = json.Unmarshal([]byte(`{
		"name": "fooBar", 
		"id": "5sryd", 
		"link_karma": 31, 
		"comment_karma": 557, 
		"is_mod": true, 
		"has_mod_mail": true,
		"has_mail": true, 
		"inbox_count": 2,
		"created": 123456789.0, 
		"created_utc": 1315269998.0, 
		"is_gold": true,
		"gold_credits": 45,
		"gold_expiration": 23273468.0, 
		"over_18": true,
		"has_verified_email": true,
		"hide_from_robots": true
	}`), actual)
	if err != nil {
		t.Fatalf("failed to unmarshal Account: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(
			"conversion of Account failed; expected %v, got %v",
			expected,
			actual)
	}
}
