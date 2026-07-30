package model

import (
	"context"
	"fmt"
	"os/exec"
	"sync"

	"github.com/cookiengineer/systempanel/parsers/pactl"
)

type VolumeModel struct {
	mu        sync.Mutex
	observers []Observer
}

type VolumeDevice struct {
	Name        string
	Description string
	Volume      int
	Mute        bool
	IsInput     bool
}

func NewVolumeModel() *VolumeModel {
	return &VolumeModel{}
}

func (m *VolumeModel) Refresh(ctx context.Context) error {
	return nil
}

func (m *VolumeModel) Observe(fn Observer) func() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.observers = append(m.observers, fn)
	return func() {}
}

func (m *VolumeModel) ListSinks() ([]VolumeDevice, error) {
	out, err := exec.Command("pactl", "list", "sinks").Output()
	if err != nil {
		return nil, err
	}
	parsed := pactl.ParseSinks(string(out))
	var devices []VolumeDevice
	for _, s := range parsed {
		devices = append(devices, VolumeDevice{
			Name:        s.Name,
			Description: s.Description,
			Volume:      s.Volume,
			Mute:        s.Mute,
			IsInput:     false,
		})
	}
	return devices, nil
}

func (m *VolumeModel) ListSources() ([]VolumeDevice, error) {
	out, err := exec.Command("pactl", "list", "sources").Output()
	if err != nil {
		return nil, err
	}
	parsed := pactl.ParseSources(string(out))
	var devices []VolumeDevice
	for _, s := range parsed {
		devices = append(devices, VolumeDevice{
			Name:        s.Name,
			Description: s.Description,
			Volume:      s.Volume,
			Mute:        s.Mute,
			IsInput:     true,
		})
	}
	return devices, nil
}

func (m *VolumeModel) SetSinkVolume(name string, volume int) error {
	return exec.Command("pactl", "set-sink-volume", name, fmt.Sprintf("%d%%", volume)).Run()
}

func (m *VolumeModel) ToggleMuteSink(name string) error {
	return exec.Command("pactl", "set-sink-mute", name, "toggle").Run()
}

func (m *VolumeModel) ToggleMuteSource(name string) error {
	return exec.Command("pactl", "set-source-mute", name, "toggle").Run()
}
