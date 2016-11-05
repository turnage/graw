package api

import (
	"github.com/turnage/graw/internal/reap"
)

// mockReaper saves the paths it is sent and returns preconfigured results.
type mockReaper struct {
	path string
	h    reap.Harvest
	err  error
}

func (m *mockReaper) Reap(path string, _ map[string]string) (reap.Harvest, error) {
	m.path = path
	return m.h, m.err
}

func (m *mockReaper) Sow(path string, _ map[string]string) error {
	m.path = path
	return m.err
}

func reaperWhich(h reap.Harvest, err error) *mockReaper {
	return &mockReaper{
		h:   h,
		err: err,
	}
}
