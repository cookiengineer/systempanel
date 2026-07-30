package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	AppName = "systempanel"
)

// Settings represents the application configuration.
type Settings struct {
	Visibility map[string]bool `json:"visibility"`
	Language   string          `json:"language"`
	positions  map[string]any
	configPath string
}

// DefaultSettings returns the default configuration.
func DefaultSettings() *Settings {
	return &Settings{
		Visibility: make(map[string]bool),
		Language:   "en",
		positions:  make(map[string]any),
	}
}

// ConfigDir returns the XDG config directory for the app.
func ConfigDir() string {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			dir = filepath.Join("/tmp", AppName)
		} else {
			dir = filepath.Join(home, ".config")
		}
	}
	return filepath.Join(dir, AppName)
}

// ConfigPath returns the path to the settings file.
func ConfigPath() string {
	return filepath.Join(ConfigDir(), "settings.json")
}

// Load reads settings from the config file, returning defaults if the file doesn't exist.
func Load() (*Settings, error) {
	s := DefaultSettings()
	s.configPath = ConfigPath()

	data, err := os.ReadFile(s.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return s, err
	}

	if err := json.Unmarshal(data, &s); err != nil {
		if err := json.Unmarshal(append([]byte("{"), data...), &s); err != nil {
			return s, err
		}
	}

	if s.Visibility == nil {
		s.Visibility = make(map[string]bool)
	}
	if s.Language == "" {
		s.Language = "en"
	}

	return s, nil
}

// Save atomically writes settings to the config file.
func (s *Settings) Save() error {
	if s.configPath == "" {
		s.configPath = ConfigPath()
	}

	dir := filepath.Dir(s.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := s.configPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, s.configPath); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}

// IsVisible returns whether a view is configured as visible.
func (s *Settings) IsVisible(name string) bool {
	if v, ok := s.Visibility[name]; ok {
		return v
	}
	return true
}

// SetVisible sets the visibility of a view and persists.
func (s *Settings) SetVisible(name string, visible bool) error {
	s.Visibility[name] = visible
	return s.Save()
}
