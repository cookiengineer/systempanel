package gtktheme

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Settings struct {
	ThemeName     string
	IconThemeName string
	PreferDark    bool
	Path          string
}

func Load(path string) (*Settings, error) {
	s := &Settings{Path: path}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			continue
		}
		eq := strings.Index(line, "=")
		if eq < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		value := strings.TrimSpace(line[eq+1:])
		if key == "gtk-theme-name" {
			s.ThemeName = value
		}
		if key == "gtk-icon-theme-name" {
			s.IconThemeName = value
		}
		if key == "gtk-application-prefer-dark-theme" {
			s.PreferDark = value == "1" || value == "true"
		}
	}
	return s, scanner.Err()
}

func Save(path, themeName string) error {
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	var content string
	data, err := os.ReadFile(path)
	if err != nil {
		content = "[Settings]\ngtk-theme-name=" + themeName + "\n"
	} else {
		content = setKey(string(data), "gtk-theme-name", themeName)
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func SaveIcon(path, iconThemeName string) error {
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	var content string
	data, err := os.ReadFile(path)
	if err != nil {
		content = "[Settings]\ngtk-icon-theme-name=" + iconThemeName + "\n"
	} else {
		content = setKey(string(data), "gtk-icon-theme-name", iconThemeName)
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func setKey(content, key, value string) string {
	lines := strings.Split(content, "\n")
	found := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "[") {
			continue
		}
		eq := strings.Index(trimmed, "=")
		if eq < 0 {
			continue
		}
		if strings.TrimSpace(trimmed[:eq]) == key {
			lines[i] = key + "=" + value
			found = true
			break
		}
	}
	if !found {
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.HasPrefix(strings.TrimSpace(lines[i]), "[") {
				lines = append(lines[:i+1], append([]string{key + "=" + value}, lines[i+1:]...)...)
				break
			}
		}
	}
	return strings.Join(lines, "\n")
}

func SaveDarkTheme(dark bool) error {
	cfg, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	val := "0"
	if dark {
		val = "1"
	}
	for _, ver := range []string{"gtk-3.0", "gtk-4.0"} {
		dir := filepath.Join(cfg, ver)
		os.MkdirAll(dir, 0755)
		path := filepath.Join(dir, "settings.ini")
		var content string
		data, err := os.ReadFile(path)
		if err != nil {
			content = "[Settings]\ngtk-application-prefer-dark-theme=" + val + "\n"
		} else {
			content = setKey(string(data), "gtk-application-prefer-dark-theme", val)
		}
		if werr := os.WriteFile(path, []byte(content), 0644); werr != nil {
			return werr
		}
	}
	return nil
}
