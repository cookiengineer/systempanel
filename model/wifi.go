package model

import (
	"context"
	"os/exec"
	"strings"

	"github.com/cookiengineer/systempanel/parsers/nmcli"
)

type WiFiNetwork struct {
	SSID     string
	BSSID    string
	Signal   int
	Security string
	Active   bool
}

type WiFiModel struct {
	observers []Observer
}

func (m *WiFiModel) Refresh(ctx context.Context) error { return nil }
func (m *WiFiModel) Observe(fn Observer) func()        { return func() {} }

func (m *WiFiModel) IsServiceRunning() bool {
	return exec.Command("systemctl", "is-active", "--quiet", "NetworkManager.service").Run() == nil
}

func (m *WiFiModel) Scan() ([]WiFiNetwork, error) {
	out, err := exec.Command("nmcli", "-t", "-f", "IN-USE,SSID,MODE,CHAN,RATE,SIGNAL,SECURITY,BSSID", "device", "wifi", "list").Output()
	if err != nil {
		return nil, err
	}
	parsed := nmcli.ParseNetworks(string(out))
	var networks []WiFiNetwork
	for _, n := range parsed {
		networks = append(networks, WiFiNetwork{
			SSID:     n.SSID,
			BSSID:    n.BSSID,
			Signal:   n.Signal,
			Security: n.Security,
			Active:   n.Active,
		})
	}
	return networks, nil
}

func (m *WiFiModel) Connect(ssid, password string) error {
	args := []string{"device", "wifi", "connect", ssid}
	if password != "" {
		args = append(args, "password", password)
	}
	return exec.Command("nmcli", args...).Run()
}

func (m *WiFiModel) Disconnect(ssid string) error {
	return exec.Command("nmcli", "connection", "down", ssid).Run()
}

func (m *WiFiModel) WiFiEnabled() bool {
	out, err := exec.Command("nmcli", "radio", "wifi").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "enabled"
}

func (m *WiFiModel) ConfiguredSSIDs() (map[string]bool, error) {
	out, err := exec.Command("nmcli", "-t", "-f", "NAME,TYPE", "connection", "show").Output()
	if err != nil {
		return nil, err
	}
	profiles := make(map[string]bool)
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		parts := strings.Split(line, ":")
		if len(parts) >= 2 && parts[1] == "802-11-wireless" {
			profiles[parts[0]] = true
		}
	}
	return profiles, nil
}
