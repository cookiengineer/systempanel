package model

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/cookiengineer/systempanel/parsers/xrandr"
)

type Monitor struct {
	Name       string
	Connected  bool
	Primary    bool
	Resolution string
	X          int
	Y          int
	Modes      []ResolutionMode
}

type ResolutionMode struct {
	Width     int
	Height    int
	Refresh   float64
	Current   bool
	Preferred bool
}

type MonitorModel struct {
	observers []Observer
}

func (m *MonitorModel) Refresh(ctx context.Context) error { return nil }
func (m *MonitorModel) Observe(fn Observer) func()        { return func() {} }

func (m *MonitorModel) ListMonitors() ([]Monitor, error) {
	out, err := exec.Command("xrandr").Output()
	if err != nil {
		return nil, err
	}
	_, outputs, err := xrandr.Parse(string(out))
	if err != nil {
		return nil, err
	}
	var monitors []Monitor
	for _, o := range outputs {
		var modes []ResolutionMode
		for _, md := range o.Modes {
			modes = append(modes, ResolutionMode{
				Width:     md.Width,
				Height:    md.Height,
				Refresh:   md.Refresh,
				Current:   md.Current,
				Preferred: md.Preferred,
			})
		}
		monitors = append(monitors, Monitor{
			Name:       o.Name,
			Connected:  o.Connected,
			Primary:    o.Primary,
			Resolution: formatRes(o.Width, o.Height),
			X:          o.X,
			Y:          o.Y,
			Modes:      modes,
		})
	}
	return monitors, nil
}

func (m *MonitorModel) SetResolution(output, mode string) error {
	return exec.Command("xrandr", "--output", output, "--mode", mode).Run()
}

func (m *MonitorModel) SetPosition(output, relation, target string) error {
	return exec.Command("xrandr", "--output", output, "--"+relation, target).Run()
}

func formatRes(w, h int) string {
	if w == 0 && h == 0 {
		return ""
	}
	return fmt.Sprintf("%dx%d", w, h)
}
