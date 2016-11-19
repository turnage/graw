package api

import (
	"github.com/turnage/graw/internal/reap"
)

// mockReaper saves the paths it is sent and returns preconfigured results.
type mockReaper struct {
	// Path is the path received by the most recent Reap or Sow call.
	Path string

	h   reap.Harvest
	err error
}

func (m *mockReaper) Reap(path string, _ map[string]string) (reap.Harvest, error) {
	m.Path = path
	return m.h, m.err
}

func (m *mockReaper) Sow(path string, _ map[string]string) error {
	m.Path = path
	return m.err
}

func ReaperWhich(h reap.Harvest, err error) *mockReaper {
	return &mockReaper{
		h:   h,
		err: err,
	}
}
