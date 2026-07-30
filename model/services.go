package model

import (
	"context"
	"os/exec"

	"github.com/cookiengineer/systempanel/parsers/systemctl"
)

type SystemdUnit struct {
	Name        string
	Description string
	LoadState   string
	ActiveState string
	SubState    string
	UnitType    string
}

type ServicesModel struct {
	observers []Observer
}

func (m *ServicesModel) Refresh(ctx context.Context) error { return nil }
func (m *ServicesModel) Observe(fn Observer) func()        { return func() {} }

func (m *ServicesModel) ListUnits() ([]SystemdUnit, error) {
	out, err := exec.Command("systemctl", "list-units", "--all", "--no-legend", "--no-pager").Output()
	if err != nil {
		return nil, err
	}
	parsed := systemctl.ParseUnits(string(out))
	var units []SystemdUnit
	for _, u := range parsed {
		units = append(units, SystemdUnit{
			Name:        u.Name,
			Description: u.Description,
			LoadState:   u.LoadState,
			ActiveState: u.ActiveState,
			SubState:    u.SubState,
			UnitType:    u.UnitType,
		})
	}
	return units, nil
}

func (m *ServicesModel) Start(name string) error {
	return exec.Command("systemctl", "start", name).Run()
}

func (m *ServicesModel) Stop(name string) error {
	return exec.Command("systemctl", "stop", name).Run()
}

func (m *ServicesModel) Restart(name string) error {
	return exec.Command("systemctl", "restart", name).Run()
}

func (m *ServicesModel) Enable(name string) error {
	return exec.Command("systemctl", "enable", name).Run()
}

func (m *ServicesModel) Disable(name string) error {
	return exec.Command("systemctl", "disble", name).Run()
}
