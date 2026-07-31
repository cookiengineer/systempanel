package model

import (
	"context"
	"os/exec"
	"time"

	"github.com/cookiengineer/systempanel/parsers/bluetoothctl"
)

type BluetoothDevice struct {
	MAC       string
	Name      string
	Paired    bool
	Connected bool
	Trusted   bool
	RSSI      int
}

type BluetoothModel struct {
	observers []Observer
	scanning  bool
	stopCh    chan struct{}
}

func (m *BluetoothModel) Refresh(ctx context.Context) error { return nil }
func (m *BluetoothModel) Observe(fn Observer) func()        { return func() {} }

func (m *BluetoothModel) IsServiceRunning() bool {
	return exec.Command("systemctl", "is-active", "--quiet", "bluetooth.service").Run() == nil
}

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
			RSSI:      info.RSSI,
		})
	}
	return devices, nil
}

func (m *BluetoothModel) powerOn() error {
	return exec.Command("bluetoothctl", "power", "on").Run()
}

func (m *BluetoothModel) setDiscoverable(on bool) error {
	if on {
		return exec.Command("bluetoothctl", "discoverable", "on").Run()
	}
	return exec.Command("bluetoothctl", "discoverable", "off").Run()
}

func (m *BluetoothModel) startScan() error {
	return exec.Command("bluetoothctl", "scan", "on").Start()
}

func (m *BluetoothModel) stopScan() error {
	return exec.Command("bluetoothctl", "scan", "off").Run()
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

func RSSIToPercent(rssi int) int {
	if rssi == 0 {
		return 0
	}
	if rssi > 0 {
		return 0
	}
	pct := ((rssi + 100) * 100) / 70
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	return pct
}

func (m *BluetoothModel) ScanLoop(onDevice func([]BluetoothDevice), onDone func()) {
	m.powerOn()
	m.setDiscoverable(true)
	m.startScan()

	pollInterval := 2 * time.Second
	scanDuration := 30 * time.Second
	deadline := time.After(scanDuration)
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			m.stopScan()
			devices, _ := m.ListDevices()
			if onDevice != nil {
				onDevice(devices)
			}
			if onDone != nil {
				onDone()
			}
			return
		case <-ticker.C:
			devices, err := m.ListDevices()
			if err == nil && onDevice != nil {
				onDevice(devices)
			}
		}
	}
}
