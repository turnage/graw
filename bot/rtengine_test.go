package bot

import (
	"container/list"
	"testing"

	"github.com/turnage/graw/bot/internal/client"
	"github.com/turnage/graw/bot/internal/operator"
)

func TestTip(t *testing.T) {
	eng := &rtEngine{
		op: operator.New(
			client.NewMock(`{
				"data": {
					"children": [
						{"data":{"name":"4"}},
						{"data":{"name":"3"}},
						{"data":{"name":"2"}},
						{"data":{"name":"1"}}
					]
				}
			}`),
		),
	}
	tips := list.New()
	tips.PushFront("shouldnotbeatback")
	for i := 0; i < fallbackCount-1; i++ {
		tips.PushFront("bunk")
	}

	posts, err := eng.tip("", tips, 0)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if tips.Back().Value.(string) == "shouldnotbeatback" {
		t.Errorf("tips were not truncated at capacity")
	}

	if tips.Front().Value.(string) != "4" {
		t.Errorf("wanted front tip '4'; got %s", tips.Front().Value.(string))
	}

	if len(posts) != 4 {
		t.Errorf("wanted 4 incrementally named posts; got %v", posts)
	}
}

func TestFixTip(t *testing.T) {
	eng := &rtEngine{
		op: operator.New(
			client.NewMock(`{
				"data": {
					"children": [
						{"data":{"name":"1"}},
						{"data":{"name":"2"}},
						{"data":{"name":"3"}},
						{"data":{"name":"4"}}
					]
				}
			}`),
		),
	}
	tips := list.New()
	tips.PushFront("1")
	tips.PushFront("0")
	broken, err := eng.fixTip(tips)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if tips.Front().Value.(string) != "1" {
		t.Errorf("wanted '1'; got %s", tips.Front().Value.(string))
	}
	if !broken {
		t.Errorf("wanted fixTip to indicate broken; 0 was gone")
	}

	tips = list.New()
	tips.PushFront("mark")
	tips.PushFront("internet")
	_, err = eng.fixTip(tips)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if tips.Front().Value.(string) != "" {
		t.Errorf("wanted ''; got %s", tips.Front().Value.(string))
	}

	tips = list.New()
	tips.PushFront("1")
	broken, err = eng.fixTip(tips)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if broken {
		t.Errorf("wanted broken to be false; '1' was valid tip")
	}
}
