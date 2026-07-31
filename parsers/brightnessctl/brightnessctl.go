package brightnessctl

import (
	"strconv"
	"strings"
)

type Device struct {
	Name       string
	Class      string
	Current    int
	Max        int
	Percentage int
}

func ParseList(output string) []Device {
	var devices []Device
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		dev := parseLine(line)
		if dev != nil {
			devices = append(devices, *dev)
		}
	}
	return devices
}

func ParseInfo(output string) *Device {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		dev := parseLine(line)
		if dev != nil {
			return dev
		}
	}
	return nil
}

func parseLine(line string) *Device {
	parts := strings.Split(line, ",")
	if len(parts) < 5 {
		return nil
	}
	current, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return nil
	}
	max, err := strconv.Atoi(strings.TrimSpace(parts[4]))
	if err != nil {
		return nil
	}
	percentStr := strings.TrimSuffix(strings.TrimSpace(parts[3]), "%")
	percentage, err := strconv.Atoi(percentStr)
	if err != nil {
		percentage = 0
		if max > 0 {
			percentage = int(float64(current) / float64(max) * 100)
		}
	}
	return &Device{
		Name:       strings.TrimSpace(parts[0]),
		Class:      strings.TrimSpace(parts[1]),
		Current:    current,
		Max:        max,
		Percentage: percentage,
	}
}
