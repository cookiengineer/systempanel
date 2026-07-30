package pactl

import (
	"bufio"
	"strings"
)

// Sink represents a PulseAudio output device (sink).
type Sink struct {
	Index       int
	Name        string
	Description string
	Mute        bool
	Volume      int
	BaseVolume  int
	Driver      string
	SampleSpec  string
	ChannelMap  string
	OwnerModule int
	MonitorSource string
	Latency     string
	Flags       string
	Properties  map[string]string
	Ports       []Port
	ActivePort  string
}

// Source represents a PulseAudio input device (source).
type Source struct {
	Index       int
	Name        string
	Description string
	Mute        bool
	Volume      int
	BaseVolume  int
	Driver      string
	SampleSpec  string
	ChannelMap  string
	OwnerModule int
	MonitorOfSink string
	Latency     string
	Flags       string
	Properties  map[string]string
	Ports       []Port
	ActivePort  string
}

// Card represents a PulseAudio card (hardware device).
type Card struct {
	Index    int
	Name     string
	Driver   string
	OwnerModule int
	Properties map[string]string
	Profiles  []Profile
	ActiveProfile string
}

// Port represents an audio port on a sink or source.
type Port struct {
	Name        string
	Description string
	Priority    int
	Available   string
}

// Profile represents an audio profile on a card.
type Profile struct {
	Name        string
	Description string
	Priority    int
	Available   string
}

// SinkInput represents an application streaming audio to a sink.
type SinkInput struct {
	Index      int
	Driver     string
	OwnerModule int
	Client     int
	Sink       int
	SampleSpec string
	ChannelMap string
	Mute       bool
	Volume     int
	ResampleMethod string
	Properties map[string]string
}

// ParseSinks parses the output of "pactl list sinks".
func ParseSinks(output string) []Sink {
	var sinks []Sink
	var current *Sink
	currentSection := ""

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "Sink #") {
			if current != nil {
				sinks = append(sinks, *current)
			}
			current = &Sink{Properties: make(map[string]string), Ports: make([]Port, 0)}
			currentSection = "sink"
			continue
		}
		if current == nil {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " \t"))

		switch {
		case indent <= 4:
			currentSection = parseSection(trimmed, current)
		case currentSection == "ports" && strings.HasPrefix(trimmed, "["):
			port := parsePort(trimmed)
			current.Ports = append(current.Ports, port)
		case strings.HasPrefix(trimmed, "Properties:") || strings.Contains(trimmed, "device.icon_name"):
			if kv := parseProperty(trimmed); kv != nil {
				current.Properties[kv[0]] = kv[1]
			}
		case strings.HasPrefix(trimmed, "Active Port:") && currentSection == "ports":
			current.ActivePort = strings.TrimSpace(strings.TrimPrefix(trimmed, "Active Port:"))
		}
	}
	if current != nil {
		sinks = append(sinks, *current)
	}
	return sinks
}

func parseSection(line string, sink *Sink) string {
	switch {
	case strings.HasPrefix(line, "Name:"):
		sink.Name = strings.TrimSpace(line[5:])
	case strings.HasPrefix(line, "Description:"):
		sink.Description = strings.TrimSpace(line[12:])
	case strings.HasPrefix(line, "Driver:"):
		sink.Driver = strings.TrimSpace(line[7:])
	case strings.HasPrefix(line, "Sample Specification:"):
		sink.SampleSpec = strings.TrimSpace(line[21:])
	case strings.HasPrefix(line, "Channel Map:"):
		sink.ChannelMap = strings.TrimSpace(line[12:])
	case strings.HasPrefix(line, "Owner Module:"):
		sink.OwnerModule = parseInt(strings.TrimSpace(line[13:]))
	case strings.HasPrefix(line, "Mute:"):
		sink.Mute = strings.TrimSpace(line[5:]) == "yes"
	case strings.HasPrefix(line, "Volume:"):
		parseVolume(line, &sink.Volume, &sink.BaseVolume)
	case strings.HasPrefix(line, "Monitor Source:"):
		sink.MonitorSource = strings.TrimSpace(line[15:])
	case strings.HasPrefix(line, "Latency:"):
		sink.Latency = strings.TrimSpace(line[8:])
	case strings.HasPrefix(line, "Flags:"):
		sink.Flags = strings.TrimSpace(line[6:])
	case strings.HasPrefix(line, "Ports:"):
		return "ports"
	}
	return ""
}

func parseVolume(line string, volume, baseVolume *int) {
	parts := strings.Fields(line)
	for _, p := range parts {
		if strings.HasSuffix(p, "%") {
			*volume = parseInt(strings.TrimSuffix(p, "%"))
		}
	}
}

func parsePort(line string) Port {
	var p Port
	line = strings.TrimSpace(line)
	if nameEnd := strings.Index(line, "]"); nameEnd > 1 {
		p.Name = line[1:nameEnd]
	}
	if descIdx := strings.Index(line, ":"); descIdx > 0 {
		p.Description = strings.TrimSpace(line[descIdx+1:])
	}
	if strings.Contains(line, "available: yes") {
		p.Available = "yes"
	} else if strings.Contains(line, "available: no") {
		p.Available = "no"
	}
	return p
}

func parseProperty(line string) []string {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return nil
	}
	key := strings.Trim(parts[0], ` "`)
	val := strings.Trim(parts[1], ` "`)
	return []string{key, val}
}

func parseInt(s string) int {
	s = strings.TrimSpace(s)
	s = strings.TrimRight(s, "%")
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
