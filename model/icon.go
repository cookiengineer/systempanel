package model

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/parsers/gtktheme"
)

type Icon struct {
	Name string
	Path string
}

type IconModel struct {
	observers []Observer
}

func (m *IconModel) Refresh(ctx context.Context) error { return nil }
func (m *IconModel) Observe(fn Observer) func()        { return func() {} }

func (m *IconModel) ListIcons() ([]Icon, error) {
	var icons []Icon
	seen := make(map[string]bool)

	dirs := []string{
		filepath.Join(xdgDataHome(), "icons"),
		filepath.Join(homeDir(), ".icons"),
		"/usr/share/icons",
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
			iconPath := filepath.Join(dir, e.Name())
			if !fileExists(filepath.Join(iconPath, "index.theme")) {
				continue
			}
			seen[e.Name()] = true
			icons = append(icons, Icon{
				Name: e.Name(),
				Path: iconPath,
			})
		}
	}
	return icons, nil
}

func (m *IconModel) CurrentIcon() string {
	path := ConfigPath()
	settings, err := gtktheme.Load(path)
	if err != nil || settings.IconThemeName == "" {
		return ""
	}
	return settings.IconThemeName
}

func (m *IconModel) ApplyIcon(name string) error {
	path := ConfigPath()
	themeName := currentThemeName()
	if err := gtktheme.SaveIcon(path, name); err != nil {
		return err
	}
	if themeName != "" {
		gtktheme.Save(path, themeName)
	}

	go func() {
		exec.Command("gsettings", "set", "org.gnome.desktop.interface", "icon-theme", name).Run()
		if themeName != "" {
			exec.Command("gsettings", "set", "org.gnome.desktop.interface", "gtk-theme", themeName).Run()
		}
	}()

	gtk4.IdleAdd(func() {
		gtk4.SetIconThemeName(name)
	})

	return nil
}
