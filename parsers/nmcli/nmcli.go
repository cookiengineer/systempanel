package nmcli

import (
	"bufio"
	"strconv"
	"strings"
)

// Network represents a Wi-Fi network scanned via nmcli.
type Network struct {
	Active   bool
	SSID     string
	BSSID    string
	Mode     string
	Channel  int
	Freq     int
	Rate     int
	Signal   int
	Security string
	Bars     string
}

// ParseNetworks parses the output of "nmcli -t -f ALL device wifi list".
// The -t flag produces terse output with : separators.
func ParseNetworks(output string) []Network {
	var networks []Network
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "IN-USE") {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue
		}
		n := Network{}
		n.Active = fields[0] == "*"
		n.SSID = unescape(fields[1])
		n.Mode = fields[2]
		n.Channel, _ = strconv.Atoi(fields[3])
		n.Rate, _ = strconv.Atoi(fields[4])
		n.Signal, _ = strconv.Atoi(fields[5])
		n.Security = strings.Join(fields[6:], " ")
		n.BSSID = ""
		if len(fields) > 7 {
			n.BSSID = fields[7]
		}
		if len(fields) > 8 {
			n.Bars = fields[8]
		}
		if n.Active {
			networks = append([]Network{n}, networks...)
		} else {
			networks = append(networks, n)
		}
	}
	return networks
}

// NetworkConnection represents a saved network connection.
type NetworkConnection struct {
	Name    string
	UUID    string
	Type    string
	Device  string
}

// ParseConnections parses "nmcli -t connection show" output.
func ParseConnections(output string) []NetworkConnection {
	var conns []NetworkConnection
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) < 4 {
			continue
		}
		conns = append(conns, NetworkConnection{
			Name:   unescape(fields[0]),
			UUID:   fields[1],
			Type:   fields[2],
			Device: fields[3],
		})
	}
	return conns
}

func unescape(s string) string {
	s = strings.ReplaceAll(s, "\\:", ":")
	s = strings.ReplaceAll(s, "\\;", ";")
	s = strings.ReplaceAll(s, "\\\\", "\\")
	return s
}
