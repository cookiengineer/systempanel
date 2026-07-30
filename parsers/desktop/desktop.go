package desktop

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

type DesktopEntry struct {
	Path  string
	Sections map[string]DesktopSection
}

type DesktopSection struct {
	Name     string
	Keys     map[string]string
	KeyOrder []string
}

func Parse(path string) (*DesktopEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	entry := &DesktopEntry{
		Path:     path,
		Sections: make(map[string]DesktopSection),
	}

	var current *DesktopSection
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			name := line[1 : len(line)-1]
			if sec, ok := entry.Sections[name]; ok {
				current = &sec
			} else {
				current = &DesktopSection{Name: name, Keys: make(map[string]string)}
				entry.Sections[name] = *current
			}
			continue
		}
		if current == nil {
			continue
		}
		eq := strings.Index(line, "=")
		if eq < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		value := strings.TrimSpace(line[eq+1:])
		sec := entry.Sections[current.Name]
		if _, exists := sec.Keys[key]; !exists {
			sec.KeyOrder = append(sec.KeyOrder, key)
		}
		sec.Keys[key] = value
		entry.Sections[current.Name] = sec
	}

	if _, ok := entry.Sections["Desktop Entry"]; !ok {
		entry.Sections["Desktop Entry"] = DesktopSection{Name: "Desktop Entry", Keys: make(map[string]string)}
	}

	return entry, scanner.Err()
}

func (e *DesktopEntry) Serialize() string {
	var b strings.Builder

	sectionOrder := []string{"Desktop Entry"}

	for _, name := range sectionOrder {
		sec, ok := e.Sections[name]
		if !ok || len(sec.Keys) == 0 {
			continue
		}
		b.WriteString("[")
		b.WriteString(name)
		b.WriteString("]\n")
		for _, key := range sec.KeyOrder {
			b.WriteString(key)
			b.WriteString("=")
			b.WriteString(sec.Keys[key])
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	for _, sec := range e.Sections {
		if sec.Name == "Desktop Entry" || len(sec.Keys) == 0 {
			continue
		}
		b.WriteString("[")
		b.WriteString(sec.Name)
		b.WriteString("]\n")
		for _, key := range sec.KeyOrder {
			b.WriteString(key)
			b.WriteString("=")
			b.WriteString(sec.Keys[key])
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (e *DesktopEntry) Get(key string) string {
	return e.SectionKey("Desktop Entry", key)
}

func (e *DesktopEntry) SectionKey(section, key string) string {
	if sec, ok := e.Sections[section]; ok {
		return sec.Keys[key]
	}
	return ""
}

func (e *DesktopEntry) Set(key, value string) {
	e.SetSectionKey("Desktop Entry", key, value)
}

func (e *DesktopEntry) SetSectionKey(section, key, value string) {
	sec, ok := e.Sections[section]
	if !ok {
		sec = DesktopSection{Name: section, Keys: make(map[string]string)}
	}
	if value == "" {
		delete(sec.Keys, key)
		sec.KeyOrder = removeKeyOrder(sec.KeyOrder, key)
	} else {
		if _, exists := sec.Keys[key]; !exists {
			sec.KeyOrder = append(sec.KeyOrder, key)
		}
		sec.Keys[key] = value
	}
	e.Sections[section] = sec
}

func (e *DesktopEntry) Keys() []string {
	return e.SectionKeys("Desktop Entry")
}

func (e *DesktopEntry) SectionKeys(section string) []string {
	if sec, ok := e.Sections[section]; ok {
		return sec.KeyOrder
	}
	return nil
}

func DesktopKnownKeys() []string {
	return []string{
		"Type", "Name", "GenericName", "NoDisplay", "Comment", "Icon",
		"Hidden", "OnlyShowIn", "NotShowIn", "DBusActivatable",
		"TryExec", "Exec", "Path", "Terminal", "Actions",
		"MimeType", "Categories", "Keywords", "StartupNotify",
		"StartupWMClass", "URL", "PrefersNonDefaultGPU",
		"SingleMainWindow",
	}
}

func SortDesktopKeys(keys []string) []string {
	known := DesktopKnownKeys()
	knownSet := make(map[string]bool)
	for _, k := range known {
		knownSet[k] = true
	}
	var knownKeys, unknown []string
	for _, k := range keys {
		if knownSet[k] {
			knownKeys = append(knownKeys, k)
		} else {
			unknown = append(unknown, k)
		}
	}
	sort.Strings(unknown)
	var result []string
	for _, k := range known {
		if containsStr(knownKeys, k) {
			result = append(result, k)
		}
	}
	result = append(result, unknown...)
	return result
}

func removeKeyOrder(keys []string, key string) []string {
	var result []string
	for _, k := range keys {
		if k != key {
			result = append(result, k)
		}
	}
	return result
}

func containsStr(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
