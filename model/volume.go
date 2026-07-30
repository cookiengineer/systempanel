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

func (m *VolumeModel) ListSinks() ([]pactl.Sink, error) {
	out, err := exec.Command("pactl", "list", "sinks").Output()
	if err != nil {
		return nil, err
	}
	return pactl.ParseSinks(string(out)), nil
}

func (m *VolumeModel) SetSinkVolume(name string, volume int) error {
	sinkName := getSinkName(name)
	return exec.Command("pactl", "set-sink-volume", sinkName, fmt.Sprintf("%d%%", volume)).Run()
}

func (m *VolumeModel) ToggleMuteSink(name string) error {
	sinkName := getSinkName(name)
	return exec.Command("pactl", "set-sink-mute", sinkName, "toggle").Run()
}

func (m *VolumeModel) SetDefaultSink(name string) error {
	sinkName := getSinkName(name)
	return exec.Command("pactl", "set-default-sink", sinkName).Run()
}

func getSinkName(name string) string {
	if len(name) > 10 {
		return name
	}
	return name
}

func (m *VolumeModel) notify() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, o := range m.observers {
		o()
	}
}
