package model

import (
	"context"
	"os/exec"
)

type PowerModel struct {
	observers []Observer
}

func (m *PowerModel) Refresh(ctx context.Context) error {
	return nil
}

func (m *PowerModel) Observe(fn Observer) func() {
	m.observers = append(m.observers, fn)
	return func() {
		for i, o := range m.observers {
			if &o == &fn {
				m.observers = append(m.observers[:i], m.observers[i+1:]...)
				return
			}
		}
	}
}

func (m *PowerModel) PowerOff() {
	exec.Command("systemctl", "poweroff").Start()
}

func (m *PowerModel) Reboot() {
	exec.Command("systemctl", "reboot").Start()
}

func (m *PowerModel) Suspend() {
	exec.Command("systemctl", "suspend").Start()
}

func (m *PowerModel) Hibernate() {
	exec.Command("systemctl", "hibernate").Start()
}
