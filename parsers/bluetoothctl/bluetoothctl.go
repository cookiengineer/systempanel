package bluetoothctl

import (
	"bufio"
	"strings"
)

// Device represents a Bluetooth device from bluetoothctl.
type Device struct {
	MAC       string
	Name      string
	Alias     string
	Class     string
	Paired    bool
	Bonded    bool
	Trusted   bool
	Blocked   bool
	Connected bool
	Legacy    bool
	UUIDs     []string
	Adapter   string
	RSSI      int
	TxPower   int
	Icon      string
}

// Controller represents a Bluetooth adapter/controller.
type Controller struct {
	MAC       string
	Name      string
	Alias     string
	Class     string
	Powered   bool
	Discoverable bool
	Pairable  bool
	Discovering  bool
	UUIDs     []string
}

// ParseDevices parses "bluetoothctl devices" output.
func ParseDevices(output string) []Device {
	var devices []Device
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 3 {
			continue
		}
		devices = append(devices, Device{
			MAC:  parts[1],
			Name: parts[2],
		})
	}
	return devices
}

// ParseInfo parses "bluetoothctl info <MAC>" output.
func ParseInfo(output string) *Device {
	dev := &Device{
		UUIDs: make([]string, 0),
	}
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		colon := strings.Index(line, ":")
		if colon < 0 {
			continue
		}
		key := strings.TrimSpace(line[:colon])
		value := strings.TrimSpace(line[colon+1:])
		switch key {
		case "Device":
			dev.MAC = value
		case "Name":
			dev.Name = value
		case "Alias":
			dev.Alias = value
		case "Class":
			dev.Class = value
		case "Paired":
			dev.Paired = value == "yes"
		case "Bonded":
			dev.Bonded = value == "yes"
		case "Trusted":
			dev.Trusted = value == "yes"
		case "Blocked":
			dev.Blocked = value == "yes"
		case "Connected":
			dev.Connected = value == "yes"
		case "LegacyPairing":
			dev.Legacy = value == "yes"
		case "UUID":
			dev.UUIDs = append(dev.UUIDs, value)
		case "RSSI":
			dev.RSSI = parseInt(value)
		case "TxPower":
			dev.TxPower = parseInt(value)
		case "Icon":
			dev.Icon = value
		}
	}
	if dev.Name == "" {
		dev.Name = dev.Alias
	}
	return dev
}

// ParseControllers parses "bluetoothctl list" output.
func ParseControllers(output string) []Controller {
	var controllers []Controller
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 3 {
			continue
		}
		controllers = append(controllers, Controller{
			MAC:  parts[1],
			Name: parts[2],
		})
	}
	return controllers
}

// ParseShow parses "bluetoothctl show" output for controller details.
func ParseShow(output string) *Controller {
	ctrl := &Controller{
		UUIDs: make([]string, 0),
	}
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		colon := strings.Index(line, ":")
		if colon < 0 {
			continue
		}
		key := strings.TrimSpace(line[:colon])
		value := strings.TrimSpace(line[colon+1:])
		switch key {
		case "Controller":
			ctrl.MAC = value
		case "Name":
			ctrl.Name = value
		case "Alias":
			ctrl.Alias = value
		case "Class":
			ctrl.Class = value
		case "Powered":
			ctrl.Powered = value == "yes"
		case "Discoverable":
			ctrl.Discoverable = value == "yes"
		case "Pairable":
			ctrl.Pairable = value == "yes"
		case "Discovering":
			ctrl.Discovering = value == "yes"
		case "UUID":
			ctrl.UUIDs = append(ctrl.UUIDs, value)
		}
	}
	return ctrl
}

func parseInt(s string) int {
	var n int
	neg := false
	for i, c := range s {
		if c == '-' && i < len(s)-1 && s[i+1] >= '0' && s[i+1] <= '9' {
			neg = true
			continue
		}
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	if neg {
		n = -n
	}
	return n
}
