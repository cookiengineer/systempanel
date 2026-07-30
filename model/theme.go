package model

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/parsers/gtktheme"
)

type Theme struct {
	Name     string
	Path     string
	IsGTK4   bool
}

type ThemeModel struct {
	observers []Observer
}

func (m *ThemeModel) Refresh(ctx context.Context) error { return nil }
func (m *ThemeModel) Observe(fn Observer) func()        { return func() {} }

func (m *ThemeModel) ListThemes() ([]Theme, error) {
	var themes []Theme
	seen := make(map[string]bool)

	dirs := []string{
		filepath.Join(xdgDataHome(), "themes"),
		filepath.Join(homeDir(), ".themes"),
		"/usr/share/themes",
	}

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() || seen[e.Name()] {
				continue
			}
			themePath := filepath.Join(dir, e.Name())
			isGTK4 := fileExists(filepath.Join(themePath, "gtk-4.0", "gtk.css"))
			isGTK3 := fileExists(filepath.Join(themePath, "gtk-3.0", "gtk.css"))
			if !isGTK4 && !isGTK3 {
				continue
			}
			seen[e.Name()] = true
			themes = append(themes, Theme{
				Name:   e.Name(),
				Path:   themePath,
				IsGTK4: isGTK4,
			})
		}
	}
	return themes, nil
}

func (m *ThemeModel) CurrentTheme() string {
	configDir := filepath.Join(xdgConfigHome(), "gtk-4.0")
	path := filepath.Join(configDir, "settings.ini")
	settings, err := gtktheme.Load(path)
	if err != nil || settings.ThemeName == "" {
		return "Adwaita"
	}
	return settings.ThemeName
}

func (m *ThemeModel) ApplyTheme(name string) error {
	path := ConfigPath()
	iconName := currentIconName()
	if err := gtktheme.Save(path, name); err != nil {
		return err
	}
	if iconName != "" {
		gtktheme.SaveIcon(path, iconName)
	}

	go func() {
		exec.Command("gsettings", "set", "org.gnome.desktop.interface", "gtk-theme", name).Run()
		if iconName != "" {
			exec.Command("gsettings", "set", "org.gnome.desktop.interface", "icon-theme", iconName).Run()
		}
	}()

	gtk4.IdleAdd(func() {
		gtk4.SetThemeName(name)
	})

	return nil
}

func xdgConfigHome() string {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		dir = filepath.Join(homeDir(), ".config")
	}
	return dir
}

func xdgDataHome() string {
	dir := os.Getenv("XDG_DATA_HOME")
	if dir == "" {
		dir = filepath.Join(homeDir(), ".local", "share")
	}
	return dir
}

func homeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ConfigDir() string {
	return filepath.Join(xdgConfigHome(), "gtk-4.0")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "settings.ini")
}

func currentIconName() string {
	settings, err := gtktheme.Load(ConfigPath())
	if err != nil || settings.IconThemeName == "" {
		return ""
	}
	return settings.IconThemeName
}

func currentThemeName() string {
	settings, err := gtktheme.Load(ConfigPath())
	if err != nil || settings.ThemeName == "" {
		return ""
	}
	return settings.ThemeName
}
