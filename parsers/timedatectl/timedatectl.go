package timedatectl

import "strings"

type Status struct {
	Timezone         string
	LocalRTC         string
	CanNTP           string
	NTP              string
	NTPSynchronized  string
	TimeUSec         string
	RTCTimeUSec      string
}

func ParseShow(output string) *Status {
	s := &Status{}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "Timezone":
			s.Timezone = val
		case "LocalRTC":
			s.LocalRTC = val
		case "CanNTP":
			s.CanNTP = val
		case "NTP":
			s.NTP = val
		case "NTPSynchronized":
			s.NTPSynchronized = val
		case "TimeUSec":
			s.TimeUSec = val
		case "RTCTimeUSec":
			s.RTCTimeUSec = val
		}
	}
	return s
}

func ParseTimezones(output string) []string {
	var zones []string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		zones = append(zones, line)
	}
	return zones
}
