package monitor

import (
	"reflect"
	"testing"

	"github.com/turnage/graw/reddit"
)

type mockScanner struct{}

func (m *mockScanner) Listing(_, _ string) (reddit.Harvest, error) {
	return reddit.Harvest{}, nil
}

type mockSorter struct {
	names []string
}

func (m *mockSorter) Sort(_ reddit.Harvest) []string { return m.names }

func TestNew(t *testing.T) {
	m, err := New(Config{Scanner: &mockScanner{}, Sorter: &mockSorter{}})
	if err != nil {
		t.Errorf("error creating monitor: %v", err)
	}
	impl := m.(*monitor)
	if len(impl.tip) < 0 || impl.tip[0] != "" {
		t.Errorf("tip wrongly initialized; got %v", impl.tip)
	}

	names := []string{"1", "2"}
	m, err = New(Config{Scanner: &mockScanner{}, Sorter: &mockSorter{names}})
	if err != nil {
		t.Errorf("error creating monitor: %v", err)
	}
	impl = m.(*monitor)
	if len(impl.tip) != len(names) || !reflect.DeepEqual(names, impl.tip) {
		t.Errorf("tip wrongly filled; got %v", impl.tip)
	}
}

func TestShaveTip(t *testing.T) {
	m := &monitor{
		blanks:         1,
		blankThreshold: 1,
		tip:            []string{"1", "2"},
		scanner:        &mockScanner{},
		sorter:         &mockSorter{},
	}

	_, err := m.Update()
	if err != nil {
		t.Errorf("error in update: %v", err)
	}

	expected := []string{"1"}
	if !reflect.DeepEqual(m.tip, expected) {
		t.Errorf("wanted tip shaved; got %v", m.tip, expected)
	}

	if m.blanks != 1 {
		t.Errorf("did not want blanks reset for bad check")
	}
}

func TestStoreTip(t *testing.T) {
	m := &monitor{
		blanks:         0,
		blankThreshold: 1,
		tip:            []string{"1", "2"},
		scanner:        &mockScanner{},
		sorter:         &mockSorter{[]string{"0"}},
	}

	_, err := m.Update()
	if err != nil {
		t.Errorf("error in update: %v", err)
	}

	expected := []string{"0", "1", "2"}
	if len(m.tip) != 3 || !reflect.DeepEqual(m.tip, expected) {
		t.Errorf("wanted tip expanded; got %v", m.tip)
	}

	if m.blanks != 0 {
		t.Errorf("wanted blanks reset; got %d", m.blanks)
	}
}

func TestBackoff(t *testing.T) {
	m := &monitor{
		blanks:         1,
		blankThreshold: 1,
		tip:            []string{"1", "2"},
		scanner:        &mockScanner{},
		sorter:         &mockSorter{names: []string{"1", "2"}},
	}

	_, err := m.Update()
	if err != nil {
		t.Errorf("error in update: %v", err)
	}

	expected := []string{"1", "2"}
	if len(m.tip) != 2 || !reflect.DeepEqual(m.tip, expected) {
		t.Errorf("wanted tip expanded; got %v", m.tip)
	}

	if m.blanks != 0 {
		t.Errorf("wanted blanks reset; got %d", m.blanks)
	}

	if m.blankThreshold != 2 {
		t.Errorf("wanted threshold scaled; got %d", m.blankThreshold)
	}
}

func TestTipFilter(t *testing.T) {
	m := &monitor{
		blanks:         2,
		blankThreshold: 2,
		tip:            []string{"1", "2", "3", "4"},
		scanner:        &mockScanner{},
		sorter:         &mockSorter{names: []string{"2", "4"}},
	}

	_, err := m.Update()
	if err != nil {
		t.Errorf("error in update: %v", err)
	}

	expected := []string{"2", "4"}
	if len(m.tip) != 2 || !reflect.DeepEqual(m.tip, expected) {
		t.Errorf("wanted tip filtered; got %v", m.tip)
	}

	if m.blanks != 0 {
		t.Errorf("wanted blanks reset; got %d", m.blanks)
	}

	if m.blankThreshold != 1 {
		t.Errorf("wanted threshold descaled; got %d", m.blankThreshold)
	}
}
