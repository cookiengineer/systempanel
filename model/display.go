package model

import (
	"context"
	"os/exec"
	"strconv"

	"github.com/cookiengineer/systempanel/parsers/brightnessctl"
)

type DisplayModel struct {
	observers []Observer
}

func (m *DisplayModel) Refresh(ctx context.Context) error { return nil }
func (m *DisplayModel) Observe(fn Observer) func()        { return func() {} }

func (m *DisplayModel) ListDevices() ([]brightnessctl.Device, error) {
	out, err := exec.Command("brightnessctl", "-l").Output()
	if err != nil {
		return nil, err
	}
	return brightnessctl.ParseList(string(out)), nil
}

func (m *DisplayModel) GetDevice(name string) (*brightnessctl.Device, error) {
	out, err := exec.Command("brightnessctl", "-d", name, "info").Output()
	if err != nil {
		return nil, err
	}
	return brightnessctl.ParseInfo(string(out)), nil
}

func (m *DisplayModel) SetBrightness(name string, value int) error {
	return exec.Command("brightnessctl", "-d", name, "set", strconv.Itoa(value)).Run()
}
