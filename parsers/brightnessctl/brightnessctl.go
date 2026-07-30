package brightnessctl

import "strings"

// Device represents a brightness-controllable device.
type Device struct {
	Name       string
	Path       string
	Current    int
	Max        int
	Percentage float64
}

// ParseList parses "brightnessctl -l" output listing all devices.
func ParseList(output string) []Device {
	var devices []Device
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "'", 3)
		if len(parts) < 2 {
			continue
		}
		name := parts[1]
		path := ""
		if idx := strings.Index(line, "Device '"); idx >= 0 {
			rest := line[idx+8+len(name)+1:]
			rest = strings.TrimSpace(rest)
			path = rest
		}
		devices = append(devices, Device{
			Name: name,
			Path: path,
		})
	}
	return devices
}

// ParseInfo parses "brightnessctl -d <device> info" output.
func ParseInfo(output string) *Device {
	dev := &Device{}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
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
			dev.Name = value
		case "Current brightness":
			dev.Current = parseInt(value)
		case "Max brightness":
			dev.Max = parseInt(value)
		}
	}
	if dev.Max > 0 {
		dev.Percentage = float64(dev.Current) / float64(dev.Max) * 100
	}
	return dev
}

func parseInt(s string) int {
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
