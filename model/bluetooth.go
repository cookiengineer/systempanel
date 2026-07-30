package model

import (
	"context"
	"os/exec"

	"github.com/cookiengineer/systempanel/parsers/bluetoothctl"
)

type BluetoothDevice struct {
	MAC       string
	Name      string
	Paired    bool
	Connected bool
	Trusted   bool
}

type BluetoothModel struct {
	observers []Observer
}

func (m *BluetoothModel) Refresh(ctx context.Context) error { return nil }
func (m *BluetoothModel) Observe(fn Observer) func()        { return func() {} }

func (m *BluetoothModel) ListDevices() ([]BluetoothDevice, error) {
	out, err := exec.Command("bluetoothctl", "devices").Output()
	if err != nil {
		return nil, err
	}
	parsed := bluetoothctl.ParseDevices(string(out))
	var devices []BluetoothDevice
	for _, d := range parsed {
		infoBytes, err := exec.Command("bluetoothctl", "info", d.MAC).Output()
		if err != nil {
			continue
		}
		info := bluetoothctl.ParseInfo(string(infoBytes))
		devices = append(devices, BluetoothDevice{
			MAC:       info.MAC,
			Name:      d.Name,
			Paired:    info.Paired,
			Connected: info.Connected,
			Trusted:   info.Trusted,
		})
	}
	return devices, nil
}

func (m *BluetoothModel) Scan() error {
	exec.Command("bluetoothctl", "scan", "on").Start()
	return nil
}

func (m *BluetoothModel) Connect(mac string) error {
	return exec.Command("bluetoothctl", "connect", mac).Run()
}

func (m *BluetoothModel) Disconnect(mac string) error {
	return exec.Command("bluetoothctl", "disconnect", mac).Run()
}

func (m *BluetoothModel) Pair(mac string) error {
	return exec.Command("bluetoothctl", "pair", mac).Run()
}

func (m *BluetoothModel) Trust(mac string) error {
	return exec.Command("bluetoothctl", "trust", mac).Run()
}

func (m *BluetoothModel) Forget(mac string) error {
	return exec.Command("bluetoothctl", "remove", mac).Run()
}
