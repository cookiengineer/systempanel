package model

import (
	"context"
	"os/exec"

	"github.com/cookiengineer/systempanel/parsers/upower"
)

type BatteryModel struct {
	observers []Observer
}

func (m *BatteryModel) Refresh(ctx context.Context) error {
	return nil
}
func (m *BatteryModel) Observe(fn Observer) func() { return func() {} }

func (m *BatteryModel) GetBatteries() ([]upower.BatteryInfo, error) {
	devices, err := exec.Command("upower", "-e").Output()
	if err != nil {
		return nil, err
	}
	var batteries []upower.BatteryInfo
	paths := upower.ParseDevices(string(devices))
	for _, path := range paths {
		out, err := exec.Command("upower", "-i", path).Output()
		if err != nil {
			continue
		}
		info := upower.ParseDeviceInfo(string(out))
		if info != nil && info.PowerSupply {
			batteries = append(batteries, *info)
		}
	}
	return batteries, nil
}
