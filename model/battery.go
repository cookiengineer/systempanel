package model

import (
	"context"
	"os/exec"
	"strings"

	"github.com/cookiengineer/systempanel/parsers/upower"
)

type BatteryInfo struct {
	Name         string
	Model        string
	Vendor       string
	Percentage   float64
	State        string
	TimeToEmpty  int64
	TimeToFull   int64
	Capacity     float64
	IsPowerSupply bool
}

type BatteryModel struct {
	observers []Observer
}

func (m *BatteryModel) Refresh(ctx context.Context) error { return nil }
func (m *BatteryModel) Observe(fn Observer) func()        { return func() {} }

func (m *BatteryModel) GetBatteries() ([]BatteryInfo, error) {
	devices, err := exec.Command("upower", "-e").Output()
	if err != nil {
		return nil, err
	}
	var batteries []BatteryInfo
	paths := upower.ParseDevices(string(devices))
	for _, path := range paths {
		out, err := exec.Command("upower", "-i", path).Output()
		if err != nil {
			continue
		}
		info := upower.ParseDeviceInfo(string(out))
		if info == nil {
			continue
		}
		if !info.PowerSupply {
			continue
		}
		if !strings.Contains(path, "battery_") {
			continue
		}
		name := extractDeviceName(path)
		if info.NativePath != "" {
			name = info.NativePath
		}
		batteries = append(batteries, BatteryInfo{
			Name:         name,
			Model:        info.Model,
			Vendor:       info.Vendor,
			Percentage:   info.Percentage,
			State:        info.State,
			TimeToEmpty:  info.TimeToEmpty,
			TimeToFull:   info.TimeToFull,
			Capacity:     info.Capacity,
			IsPowerSupply: info.PowerSupply,
		})
	}
	return batteries, nil
}

func extractDeviceName(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx < 0 {
		return path
	}
	name := path[idx+1:]
	name = strings.TrimPrefix(name, "battery_")
	return name
}
