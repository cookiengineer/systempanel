package model

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cookiengineer/systempanel/parsers/desktop"
)

type AutostartEntry struct {
	Path    string
	Name    string
	Comment string
	Exec    string
	Icon    string
	Hidden  bool
}

type AutostartModel struct {
	observers []Observer
}

func (m *AutostartModel) Refresh(ctx context.Context) error { return nil }
func (m *AutostartModel) Observe(fn Observer) func()        { return func() {} }

func (m *AutostartModel) ListEntries() ([]AutostartEntry, error) {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".config", "autostart")
	os.MkdirAll(dir, 0755)

	files, err := filepath.Glob(filepath.Join(dir, "*.desktop"))
	if err != nil {
		return nil, err
	}
	var entries []AutostartEntry
	for _, f := range files {
		entry, err := desktop.Parse(f)
		if err != nil {
			continue
		}
		name := entry.Get("Name")
		if name == "" {
			name = filepath.Base(f)
			name = name[:len(name)-len(".desktop")]
		}
		entries = append(entries, AutostartEntry{
			Path:    f,
			Name:    name,
			Comment: entry.Get("Comment"),
			Exec:    entry.Get("Exec"),
			Icon:    entry.Get("Icon"),
			Hidden:  entry.Get("Hidden") == "true",
		})
	}
	return entries, nil
}
