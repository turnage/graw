package api

import (
	"fmt"
	"net/http"
	"testing"
)

var (
	linkListingJSON = []byte(`{
		"kind": "Listing",
		"data": {
			"children": [
				{
					"kind": "t3",
					"data": {
						"title": "hello",
						"body": "hello"
					}
				}
			]
		}
	}`)
	comboListingJSON = []byte(`{
		"kind": "Listing",
		"data": {
			"children": [
				{
					"kind": "t3",
					"data": {
						"title": "hello",
						"body": "hello"
					}
				},
				{
					"kind": "t1",
					"data": {
						"body": "hello"
					}
				},
				{
					"kind": "t4",
					"data": {
						"body": "hello"
					}
				}
			]
		}
	}`)
	threadJSON = []byte(`[
		{
			"kind": "Listing",
			"data": {
				"children": [
					{
						"kind": "t3",
						"data": {
							"title": "hola"
						}
					}
				]
			}
		},
		{
			"kind": "Listing",
			"data": {
				"children": [
					{
						"kind": "t1",
						"data": {
							"id": "arnold"
						}
					},
					{
						"kind": "t1",
						"data": {
							"id": "harold"
						}
					}
				]
			}
		}
	]`)
	inboxJSON = []byte(`{
		"kind": "Listing",
		"data": {
			"children" : [
				{
					"kind": "t4",
					"data": {
						"was_comment": true
					}
				}
			]
		}
	}`)
	errRequester = func(r *http.Request) ([]byte, error) {
		return nil, fmt.Errorf("error")
	}
)

func TestScrape(t *testing.T) {
	if _, _, _, err := Scrape(errRequester, "path", "", 1); err == nil {
		t.Errorf("wanted error for request error")
	}

	if posts, comments, messages, err := Scrape(
		func(r *http.Request) ([]byte, error) {
			return comboListingJSON, nil
		},
		"path", "", 1,
	); err != nil {
		t.Fatal(err)
	} else if len(posts) != 1 {
		t.Errorf("got %d posts; wanted 1", len(posts))
	} else if len(comments) != 1 {
		t.Errorf("got %d comments; wanted 1", len(comments))
	} else if len(messages) != 1 {
		t.Errorf("got %d messages; wanted 1", len(messages))
	}
}

func TestIsThereThing(t *testing.T) {
	if _, err := IsThereThing(errRequester, "1"); err == nil {
		t.Errorf("wanted error for request failure")
	}

	if exists, err := IsThereThing(
		func(r *http.Request) ([]byte, error) {
			return linkListingJSON, nil
		},
		"1",
	); err != nil {
		t.Fatalf("error: %v", err)
	} else if exists == false {
		t.Errorf("got false; wanted true")
	}

	// Malformed JSON responses should bubble an error up.
	if _, err := IsThereThing(
		func(r *http.Request) ([]byte, error) {
			return []byte(`{
				"data": {
					"children": [
						{"data": {"name": "charlie"}},
					]
				}
			}`), nil
		},
		"1",
	); err == nil {
		t.Errorf("wanted an error for a bad response")
	}

	// Missing Things should return nil.
	if exists, err := IsThereThing(
		func(r *http.Request) ([]byte, error) {
			return []byte(`{
				"kind": "Listing",
				"data": {
					"children": []
				}
			}`), nil
		},
		"1",
	); err != nil {
		t.Fatalf("error: %v", err)
	} else if exists != false {
		t.Errorf("got true; wanted false")
	}
}

func TestThread(t *testing.T) {
	if _, err := Thread(errRequester, "/thread"); err == nil {
		t.Errorf("wanted error for request failure")
	}

	thread, err := Thread(
		func(r *http.Request) ([]byte, error) {
			return threadJSON, nil
		},
		"/thread")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(thread.GetComments()) != 2 {
		t.Errorf("got %d comments; wanted 2", len(thread.GetComments()))
	}
}

func TestInbox(t *testing.T) {
	if _, err := Inbox(errRequester); err == nil {
		t.Errorf("wanted error for request failure")
	}

	if messages, err := Inbox(
		func(r *http.Request) ([]byte, error) {
			return inboxJSON, nil
		},
	); err != nil {
		t.Fatalf("error: %v", err)
	} else if len(messages) != 1 {
		t.Fatalf("got %d messages; wanted 1", len(messages))
	} else if !messages[0].GetWasComment() {
		t.Fatal("got non-comment inboxable; wanted comment inboxable")
	}
}

func TestReply(t *testing.T) {
	if err := Reply(errRequester, "parent", "content"); err == nil {
		t.Errorf("wanted error for request failure")
	}

	if err := Reply(
		func(r *http.Request) ([]byte, error) {
			return []byte(""), nil
		},
		"parent", "content",
	); err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestCompose(t *testing.T) {
	if err := Compose(errRequester, "user", "subject", "body"); err == nil {
		t.Errorf("wanted error for request failure")
	}

	if err := Compose(
		func(r *http.Request) ([]byte, error) {
			return []byte(""), nil
		},
		"user", "subject", "body",
	); err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestSubmit(t *testing.T) {
	if err := Submit(errRequester, "aww", "self", "title", ""); err == nil {
		t.Errorf("wanted error for request failure")
	}

	if err := Submit(
		func(r *http.Request) ([]byte, error) {
			return []byte(""), nil
		},
		"aww", "self", "title", "",
	); err != nil {
		t.Fatalf("error: %v", err)
	}
}
