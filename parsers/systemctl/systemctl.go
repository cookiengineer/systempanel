package systemctl

import (
	"bufio"
	"strings"
)

// Unit represents a systemd unit from "systemctl list-units --all".
type Unit struct {
	Name        string
	LoadState   string
	ActiveState string
	SubState    string
	Description string
	UnitType    string
}

// ParseUnits parses "systemctl list-units --all --no-legend" output.
func ParseUnits(output string) []Unit {
	var units []Unit
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "UNIT") {
			continue
		}
		unit := parseUnitLine(line)
		if unit != nil {
			units = append(units, *unit)
		}
	}
	return units
}

func parseUnitLine(line string) *Unit {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return nil
	}
	nameIdx := 0
	if fields[0] == "●" || fields[0] == "*" || fields[0] == " " {
		nameIdx = 1
	}
	if nameIdx >= len(fields) {
		return nil
	}
	u := &Unit{
		Name:        fields[nameIdx],
		LoadState:   "",
		ActiveState: "",
		SubState:    "",
	}
	if nameIdx+1 < len(fields) {
		u.LoadState = fields[nameIdx+1]
	}
	if nameIdx+2 < len(fields) {
		u.ActiveState = fields[nameIdx+2]
	}
	if nameIdx+3 < len(fields) {
		u.SubState = fields[nameIdx+3]
	}
	descStart := -1
	for i := range fields {
		if i > nameIdx+3 && len(fields[i]) > 0 {
			descStart = i
			break
		}
	}
	if descStart >= 0 {
		u.Description = strings.Join(fields[descStart:], " ")
	}

	dotIdx := strings.LastIndex(u.Name, ".")
	if dotIdx > 0 {
		u.UnitType = u.Name[dotIdx+1:]
	}
	return u
}

// ParseStatus parses "systemctl status <unit>" output.
func ParseStatus(output string) string {
	return strings.TrimSpace(output)
}
