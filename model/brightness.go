package model

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/cookiengineer/systempanel/parsers/brightnessctl"
)

type BrightnessDevice struct {
	Name       string
	Current    int
	Max        int
	Percentage int
}

type BrightnessModel struct {
	observers []Observer
}

func NewBrightnessModel() *BrightnessModel {
	return &BrightnessModel{}
}

func (m *BrightnessModel) Refresh(ctx context.Context) error {
	return nil
}

func (m *BrightnessModel) Observe(fn Observer) func() {
	m.observers = append(m.observers, fn)
	return func() {}
}

func (m *BrightnessModel) ListBacklights() ([]BrightnessDevice, error) {
	out, err := exec.Command("brightnessctl", "-m", "-l").Output()
	if err != nil {
		return nil, err
	}
	parsed := brightnessctl.ParseList(string(out))
	var devices []BrightnessDevice
	for _, d := range parsed {
		if d.Class != "backlight" {
			continue
		}
		devices = append(devices, BrightnessDevice{
			Name:       d.Name,
			Current:    d.Current,
			Max:        d.Max,
			Percentage: d.Percentage,
		})
	}
	return devices, nil
}

func (m *BrightnessModel) SetBrightness(name string, percentage int) error {
	return exec.Command("brightnessctl", "-d", name, "s", fmt.Sprintf("%d%%", percentage)).Run()
}
