package model

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Wallpaper struct {
	Path     string
	Filename string
}

type WallpaperModel struct {
	observers []Observer
}

func NewWallpaperModel() *WallpaperModel {
	return &WallpaperModel{}
}

func (m *WallpaperModel) Observe(fn Observer) func() {
	m.observers = append(m.observers, fn)
	return func() {
		for i, o := range m.observers {
			if &o == &fn {
				m.observers = append(m.observers[:i], m.observers[i+1:]...)
				return
			}
		}
	}
}

func (m *WallpaperModel) Refresh(ctx interface{}) error {
	return nil
}

func WallpaperDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "Pictures", "Wallpapers")
}

func IsImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".bmp", ".gif", ".webp":
		return true
	}
	return false
}

func (m *WallpaperModel) ListWallpapers() ([]Wallpaper, error) {
	dir := WallpaperDir()
	if dir == "" {
		return nil, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil
	}

	var wallpapers []Wallpaper
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if IsImageFile(e.Name()) {
			wallpapers = append(wallpapers, Wallpaper{
				Path:     filepath.Join(dir, e.Name()),
				Filename: e.Name(),
			})
		}
	}
	return wallpapers, nil
}

func (m *WallpaperModel) SetWallpaper(path string, mode string) error {
	if mode == "" {
		mode = "scale"
	}

	if _, err := os.Stat("/usr/bin/feh"); err == nil {
		return m.setWithFeh(path, mode)
	}
	return m.setWithHsetroot(path, mode)
}

func (m *WallpaperModel) setWithFeh(path string, mode string) error {
	args := []string{}
	switch mode {
	case "stretch":
		args = []string{"--bg-scale", path}
	case "scale":
		args = []string{"--bg-fill", path}
	case "center":
		args = []string{"--bg-center", path}
	case "tile":
		args = []string{"--bg-tile", path}
	case "max":
		args = []string{"--bg-max", path}
	default:
		args = []string{"--bg-fill", path}
	}
	return exec.Command("feh", args...).Run()
}

func (m *WallpaperModel) setWithHsetroot(path string, mode string) error {
	args := []string{}
	switch mode {
	case "stretch":
		args = append(args, "-stretch", path)
	case "scale":
		args = append(args, "-fill", path)
	case "center":
		args = append(args, "-center", path)
	case "tile":
		args = append(args, "-tile", path)
	case "max":
		args = append(args, "-max", path)
	default:
		args = append(args, "-fill", path)
	}
	return exec.Command("hsetroot", args...).Run()
}

func (m *WallpaperModel) GetModes() []string {
	return []string{"scale", "stretch", "center", "tile", "max"}
}

func (m *WallpaperModel) HasWallpapers() bool {
	dir := WallpaperDir()
	if dir == "" {
		return false
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return false
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && IsImageFile(e.Name()) {
			return true
		}
	}
	return false
}

func (m *WallpaperModel) HasWallpaperDir() bool {
	dir := WallpaperDir()
	if dir == "" {
		return false
	}
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return info.IsDir()
}
