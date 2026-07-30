package desktop

import (
	"bufio"
	"os"
	"strings"
)

// DesktopEntry represents a parsed XDG .desktop file.
type DesktopEntry struct {
	Path     string
	Name     string
	Comment  string
	Exec     string
	Icon     string
	Type     string
	Hidden   bool
	NoDisplay bool
	StartupNotify bool
	Terminal bool
	Categories []string
	OnlyShowIn []string
	NotShowIn  []string
}

// Parse reads and parses a .desktop file at the given path.
func Parse(path string) (*DesktopEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	entry := &DesktopEntry{Path: path}
	scanner := bufio.NewScanner(f)
	var currentSection string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line
			continue
		}
		if currentSection != "[Desktop Entry]" {
			continue
		}
		eq := strings.Index(line, "=")
		if eq < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		value := strings.TrimSpace(line[eq+1:])
		switch key {
		case "Name":
			entry.Name = value
		case "Comment":
			entry.Comment = value
		case "Exec":
			entry.Exec = value
		case "Icon":
			entry.Icon = value
		case "Type":
			entry.Type = value
		case "Hidden":
			entry.Hidden = value == "true"
		case "NoDisplay":
			entry.NoDisplay = value == "true"
		case "StartupNotify":
			entry.StartupNotify = value == "true"
		case "Terminal":
			entry.Terminal = value == "true"
		case "Categories":
			entry.Categories = strings.Split(value, ";")
		case "OnlyShowIn":
			entry.OnlyShowIn = strings.Split(value, ";")
		case "NotShowIn":
			entry.NotShowIn = strings.Split(value, ";")
		}
	}
	if entry.Name == "" {
		name := entry.Path
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			name = name[idx+1:]
		}
		name = strings.TrimSuffix(name, ".desktop")
		entry.Name = name
	}
	return entry, scanner.Err()
}
