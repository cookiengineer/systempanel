package model

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cookiengineer/systempanel/parsers/systemctl"
	"github.com/cookiengineer/systempanel/parsers/systemdunit"
)

type ServiceUnit struct {
	Name        string
	Description string
	LoadState   string
	ActiveState string
	SubState    string
	UnitType    string
	IsUser      bool
	UnitFile    *systemdunit.UnitFile
}

type ServicesModel struct {
	observers []Observer
}

func (m *ServicesModel) Refresh(ctx context.Context) error { return nil }
func (m *ServicesModel) Observe(fn Observer) func()        { return func() {} }

func (m *ServicesModel) ListUnits(user bool) ([]ServiceUnit, error) {
	if user {
		return m.listUserServices()
	}
	args := []string{"list-units", "--type=service", "--all", "--no-legend", "--no-pager"}
	out, err := exec.Command("systemctl", args...).Output()
	if err != nil {
		return nil, err
	}
	parsed := systemctl.ParseUnits(string(out))
	var units []ServiceUnit
	for _, u := range parsed {
		units = append(units, ServiceUnit{
			Name:        u.Name,
			Description: u.Description,
			LoadState:   u.LoadState,
			ActiveState: u.ActiveState,
			SubState:    u.SubState,
			UnitType:    u.UnitType,
			IsUser:      false,
		})
	}
	return units, nil
}

func (m *ServicesModel) listUserServices() ([]ServiceUnit, error) {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".config", "systemd", "user")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil
	}

	var units []ServiceUnit
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".service") || e.IsDir() {
			continue
		}
		name := e.Name()
		state := m.getUnitState(name, true)
		units = append(units, ServiceUnit{
			Name:        name,
			Description: "",
			LoadState:   state.LoadState,
			ActiveState: state.ActiveState,
			SubState:    state.SubState,
			UnitType:    "service",
			IsUser:      true,
		})
	}
	return units, nil
}

func (m *ServicesModel) getUnitState(name string, isUser bool) ServiceUnit {
	args := []string{"show", "-p", "LoadState,ActiveState,SubState,Description", name}
	if isUser {
		args = append([]string{"--user"}, args...)
	}
	out, err := exec.Command("systemctl", args...).Output()
	if err != nil {
		return ServiceUnit{}
	}
	var su ServiceUnit
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "LoadState="):
			su.LoadState = strings.TrimPrefix(line, "LoadState=")
		case strings.HasPrefix(line, "ActiveState="):
			su.ActiveState = strings.TrimPrefix(line, "ActiveState=")
		case strings.HasPrefix(line, "SubState="):
			su.SubState = strings.TrimPrefix(line, "SubState=")
		case strings.HasPrefix(line, "Description="):
			su.Description = strings.TrimPrefix(line, "Description=")
		}
	}
	return su
}

func (m *ServicesModel) LoadUnitFile(name string, isUser bool) (*systemdunit.UnitFile, error) {
	path := m.unitFilePath(name, isUser)
	return systemdunit.Parse(path)
}

func (m *ServicesModel) SaveUnitFile(name string, isUser bool, uf *systemdunit.UnitFile) error {
	content := uf.Serialize()
	path := m.unitFilePath(name, isUser)

	if isUser {
		dir := filepath.Dir(path)
		os.MkdirAll(dir, 0755)
		return os.WriteFile(path, []byte(content), 0644)
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func (m *ServicesModel) unitFilePath(name string, isUser bool) string {
	if isUser {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config", "systemd", "user", name)
	}

	if strings.HasPrefix(name, "/") {
		return name
	}

	out, err := exec.Command("systemctl", "show", "-p", "FragmentPath", name).Output()
	if err == nil {
		line := strings.TrimSpace(string(out))
		if path := strings.TrimPrefix(line, "FragmentPath="); path != "" {
			return path
		}
	}

	return filepath.Join("/etc/systemd/system", name)
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

func (m *ServicesModel) Enable(name string, isUser bool) error {
	args := []string{"enable", name}
	if isUser {
		args = append([]string{"--user"}, args...)
	}
	return exec.Command("systemctl", args...).Run()
}

func (m *ServicesModel) Disable(name string, isUser bool) error {
	args := []string{"disable", name}
	if isUser {
		args = append([]string{"--user"}, args...)
	}
	return exec.Command("systemctl", args...).Run()
}
