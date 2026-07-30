package xrandr

import (
	"bufio"
	"strconv"
	"strings"
)

type Screen struct {
	MinWidth     int
	MinHeight    int
	CurrentWidth  int
	CurrentHeight int
	MaxWidth     int
	MaxHeight    int
}

type Output struct {
	Name      string
	Connected bool
	Primary   bool
	Width     int
	Height    int
	X         int
	Y         int
	MMWidth   int
	MMHeight  int
	Modes     []Mode
}

type Mode struct {
	Width    int
	Height   int
	Refresh  float64
	Current  bool
	Preferred bool
}

func Parse(output string) (*Screen, []Output, error) {
	var screen Screen
	var outputs []Output
	var currentOutput *Output

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")

		if strings.HasPrefix(line, "Screen ") {
			screen = parseScreen(line)
			continue
		}

		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			if currentOutput != nil && currentOutput.Connected {
				mode := parseMode(strings.TrimSpace(line))
				if mode.Width > 0 {
					currentOutput.Modes = append(currentOutput.Modes, mode)
				}
			}
			continue
		}

		if strings.Contains(line, " connected") || strings.Contains(line, " disconnected") {
			output := parseOutput(line)
			if output.Name != "" {
				if currentOutput != nil {
					outputs = append(outputs, *currentOutput)
				}
				currentOutput = &output
			}
		}
	}

	if currentOutput != nil {
		outputs = append(outputs, *currentOutput)
	}

	return &screen, outputs, scanner.Err()
}

func parseScreen(line string) Screen {
	var s Screen
	rest := strings.TrimPrefix(line, "Screen 0: ")
	parts := strings.Split(rest, ", ")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "minimum ") {
			w, h := parseDims(strings.TrimPrefix(p, "minimum "))
			s.MinWidth, s.MinHeight = w, h
		} else if strings.HasPrefix(p, "current ") {
			w, h := parseDims(strings.TrimPrefix(p, "current "))
			s.CurrentWidth, s.CurrentHeight = w, h
		} else if strings.HasPrefix(p, "maximum ") {
			w, h := parseDims(strings.TrimPrefix(p, "maximum "))
			s.MaxWidth, s.MaxHeight = w, h
		}
	}
	return s
}

func parseOutput(line string) Output {
	var o Output
	if strings.Contains(line, " disconnected") {
		parts := strings.SplitN(line, " ", 2)
		o.Name = parts[0]
		o.Connected = false
		return o
	}

	o.Connected = true

	parts := strings.Fields(line)
	if len(parts) < 3 {
		return o
	}

	o.Name = parts[0]

	idx := 1
	if parts[1] == "connected" {
		idx = 2
		if parts[2] == "primary" {
			o.Primary = true
			idx = 3
		}
	}

	if idx < len(parts) {
		geo := parts[idx]
		if plusIdx := strings.Index(geo, "+"); plusIdx >= 0 {
			w, h := parseDims(geo[:plusIdx])
			o.Width, o.Height = w, h
			rest := geo[plusIdx:]
			xyParts := strings.Split(rest, "+")
			if len(xyParts) >= 3 {
				o.X, _ = strconv.Atoi(xyParts[1])
				o.Y, _ = strconv.Atoi(xyParts[2])
			}
		}
	}

	for _, p := range parts {
		if strings.HasSuffix(p, "mm") {
			mmParts := strings.SplitN(p, "x", 2)
			if len(mmParts) == 2 {
				o.MMWidth, _ = strconv.Atoi(mmParts[0])
				o.MMHeight, _ = strconv.Atoi(strings.TrimSuffix(mmParts[1], "mm"))
			}
		}
	}

	return o
}

func parseMode(line string) Mode {
	var m Mode
	parts := strings.Fields(line)
	if len(parts) < 1 {
		return m
	}

	w, h := parseDims(parts[0])
	m.Width, m.Height = w, h

	for _, p := range parts[1:] {
		freq := strings.TrimRight(p, "*+ ")
		if f, err := strconv.ParseFloat(freq, 64); err == nil {
			m.Refresh = f
		}
		if strings.Contains(p, "*") {
			m.Current = true
		}
		if strings.Contains(p, "+") {
			m.Preferred = true
		}
	}

	return m
}

func parseDims(s string) (int, int) {
	parts := strings.SplitN(s, "x", 2)
	if len(parts) != 2 {
		return 0, 0
	}
	w, _ := strconv.Atoi(parts[0])
	h, _ := strconv.Atoi(parts[1])
	return w, h
}
