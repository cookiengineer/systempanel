package model

import (
	"os/exec"
	"strings"

	"github.com/cookiengineer/systempanel/parsers/powerprofile"
)

type PowerProfile struct {
	Name     string
	Active   bool
	Driver   string
	Degraded string
}

type PowerProfileModel struct {
	observers []Observer
}

func (m *PowerProfileModel) Refresh( /* TODO */ ) error { return nil }
func (m *PowerProfileModel) Observe(fn Observer) func()  { return func() {} }

func (m *PowerProfileModel) ListProfiles() ([]PowerProfile, error) {
	out, err := exec.Command("powerprofilesctl", "list").Output()
	if err != nil {
		return nil, err
	}

	parsed := powerprofile.ParseList(string(out))
	profiles := make([]PowerProfile, len(parsed))
	for i, p := range parsed {
		profiles[i] = PowerProfile{
			Name:     p.Name,
			Active:   p.Active,
			Driver:   p.Driver,
			Degraded: p.Degraded,
		}
	}
	return profiles, nil
}

func (m *PowerProfileModel) CurrentProfile() string {
	out, err := exec.Command("powerprofilesctl", "get").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (m *PowerProfileModel) SetProfile(name string) error {
	return exec.Command("powerprofilesctl", "set", name).Run()
}
