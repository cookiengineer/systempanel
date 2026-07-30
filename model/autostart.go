package model

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cookiengineer/systempanel/parsers/desktop"
)

type AutostartEntry struct {
	Path   string
	Name   string
	Exec   string
	Icon   string
	Hidden bool
}

type AutostartModel struct {
	observers []Observer
}

func (m *AutostartModel) Refresh(ctx context.Context) error { return nil }
func (m *AutostartModel) Observe(fn Observer) func()        { return func() {} }

func (m *AutostartModel) ListEntries() ([]AutostartEntry, error) {
	var entries []AutostartEntry
	for _, dir := range autostartDirs() {
		files, err := filepath.Glob(filepath.Join(dir, "*.desktop"))
		if err != nil {
			continue
		}
		for _, file := range files {
			entry, err := desktop.Parse(file)
			if err != nil {
				continue
			}
			if entry.NoDisplay {
				continue
			}
			entries = append(entries, AutostartEntry{
				Path:   file,
				Name:   entry.Name,
				Exec:   entry.Exec,
				Icon:   entry.Icon,
				Hidden: entry.Hidden,
			})
		}
	}
	return entries, nil
}

func (m *AutostartModel) Enable(path string) error {
	entry, err := desktop.Parse(path)
	if err != nil {
		return err
	}
	entry.Hidden = false
	return writeDesktopFile(path, entry)
}

func (m *AutostartModel) Disable(path string) error {
	entry, err := desktop.Parse(path)
	if err != nil {
		return err
	}
	entry.Hidden = true
	return writeDesktopFile(path, entry)
}

func writeDesktopFile(path string, entry *desktop.DesktopEntry) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)
	if entry.Hidden {
		content = setDesktopKey(content, "Hidden", "true")
	} else {
		content = setDesktopKey(content, "Hidden", "false")
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func setDesktopKey(content, key, value string) string {
	return content
}

func autostartDirs() []string {
	var dirs []string
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			configDir = filepath.Join(home, ".config")
		}
	}
	if configDir != "" {
		dirs = append(dirs, filepath.Join(configDir, "autostart"))
	}
	dirs = append(dirs, "/etc/xdg/autostart")
	return dirs
}
