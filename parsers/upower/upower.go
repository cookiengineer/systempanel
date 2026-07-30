package upower

import (
	"bufio"
	"math"
	"strings"
)

// Device represents a UPower device (battery, line power, etc.).
type Device struct {
	Type          string
	Name          string
	Path          string
	NativePath    string
	Vendor        string
	Model         string
	Serial        string
	PowerSupply   bool
	UpdateTime    int64
	HasHistory    bool
	HasStatistics bool
	Online        bool
	IconName      string
}

// BatteryInfo represents detailed battery information from UPower.
type BatteryInfo struct {
	Device
	Percentage    float64
	State         string
	TimeToEmpty   int64
	TimeToFull    int64
	Energy        float64
	EnergyEmpty   float64
	EnergyFull    float64
	EnergyFullDesign float64
	EnergyRate    float64
	Voltage       float64
	Capacity      float64
	Technology    string
	WarningLevel  string
	ChargeCycles  int
}

// ParseDevices parses "upower -e" output (list of device paths).
func ParseDevices(output string) []string {
	var paths []string
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && strings.HasPrefix(line, "/org/freedesktop/UPower/devices/") {
			paths = append(paths, line)
		}
	}
	return paths
}

// ParseDeviceInfo parses "upower -i <path>" output for a single device.
func ParseDeviceInfo(output string) *BatteryInfo {
	info := &BatteryInfo{}
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
		case "native-path":
			info.NativePath = value
		case "vendor":
			info.Vendor = value
		case "model":
			info.Model = value
		case "serial":
			info.Serial = value
		case "power supply":
			info.PowerSupply = value == "yes"
		case "has history":
			info.HasHistory = value == "yes"
		case "has statistics":
			info.HasStatistics = value == "yes"
		case "online":
			info.Online = value == "yes"
		case "icon-name":
			info.IconName = value
		case "percentage":
			info.Percentage = parseUpowerFloat(value)
		case "state":
			info.State = value
		case "time to empty":
			info.TimeToEmpty = parseUpowerInt64(value)
		case "time to full":
			info.TimeToFull = parseUpowerInt64(value)
		case "energy":
			info.Energy = parseUpowerFloat(value)
		case "energy-empty":
			info.EnergyEmpty = parseUpowerFloat(value)
		case "energy-full":
			info.EnergyFull = parseUpowerFloat(value)
		case "energy-full-design":
			info.EnergyFullDesign = parseUpowerFloat(value)
		case "energy-rate":
			info.EnergyRate = parseUpowerFloat(value)
		case "voltage":
			info.Voltage = parseUpowerFloat(value)
		case "capacity":
			info.Capacity = parseUpowerFloat(value)
		case "technology":
			info.Technology = value
		case "warning-level":
			info.WarningLevel = value
		case "charge-cycles":
			info.ChargeCycles = int(parseUpowerFloat(value))
		}
	}
	if info.Percentage == 0 && info.EnergyFull > 0 {
		info.Percentage = math.Round(info.Energy/info.EnergyFull*10000) / 100
	}
	return info
}

// ParseAllDevices returns all UPower device paths.
func ParseAllDevices(output string) []string {
	return strings.Split(strings.TrimSpace(output), "\n")
}

func parseUpowerFloat(s string) float64 {
	var n float64
	div := 1.0
	inDecimal := false
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			if inDecimal {
				div *= 10
				n += float64(c-'0') / div
			} else {
				n = n*10 + float64(c-'0')
			}
		case c == '.':
			inDecimal = true
		}
	}
	return n
}

func parseUpowerInt64(s string) int64 {
	var n int64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int64(c-'0')
		}
	}
	return n
}
