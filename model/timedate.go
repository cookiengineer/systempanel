package model

import (
	"os/exec"

	"github.com/cookiengineer/systempanel/parsers/timedatectl"
	"github.com/cookiengineer/systempanel/widget"
)

type TimeDateStatus struct {
	Timezone        string
	NTPEnabled      bool
	NTPSynchronized bool
}

type TimeDateModel struct {
	observers []Observer
}

func (m *TimeDateModel) Refresh( /* TODO */ ) error { return nil }
func (m *TimeDateModel) Observe(fn Observer) func()  { return func() {} }

func (m *TimeDateModel) GetStatus() (*TimeDateStatus, error) {
	out, err := exec.Command("timedatectl", "show").Output()
	if err != nil {
		return nil, err
	}
	s := timedatectl.ParseShow(string(out))
	return &TimeDateStatus{
		Timezone:        s.Timezone,
		NTPEnabled:      s.NTP == "yes",
		NTPSynchronized: s.NTPSynchronized == "yes",
	}, nil
}

func (m *TimeDateModel) ListTimezones() ([]string, error) {
	out, err := exec.Command("timedatectl", "list-timezones").Output()
	if err != nil {
		return nil, err
	}
	return timedatectl.ParseTimezones(string(out)), nil
}

func (m *TimeDateModel) SetTimezone(zone string) error {
	return widget.RunSudoCommand("timedatectl", "set-timezone", zone)
}

func (m *TimeDateModel) CurrentTimezone() string {
	out, err := exec.Command("timedatectl", "show", "--property=Timezone").Output()
	if err != nil {
		return ""
	}
	s := timedatectl.ParseShow(string(out))
	return s.Timezone
}

func (m *TimeDateModel) EnableNTP() error {
	return widget.RunSudoCommand("timedatectl", "set-ntp", "true")
}

func (m *TimeDateModel) DisableNTP() error {
	return widget.RunSudoCommand("timedatectl", "set-ntp", "false")
}

func (m *TimeDateModel) SetTime(datetime string) error {
	return widget.RunSudoCommand("timedatectl", "set-time", datetime)
}
