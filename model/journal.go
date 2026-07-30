package model

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/cookiengineer/systempanel/parsers/journalctl"
)

type JournalEntry struct {
	Timestamp    time.Time
	PriorityName string
	Message      string
	Unit         string
	Hostname     string
}

type JournalModel struct {
	observers []Observer
}

func (m *JournalModel) Refresh(ctx context.Context) error { return nil }
func (m *JournalModel) Observe(fn Observer) func()        { return func() {} }

func (m *JournalModel) Fetch(count int) ([]JournalEntry, error) {
	out, err := exec.Command(
		"journalctl",
		"--no-pager",
		"--output=json",
		"-n", strconv.Itoa(int(count)),
	).Output()
	if err != nil {
		return nil, err
	}
	parsed, err := journalctl.Parse(strings.NewReader(string(out)))
	if err != nil {
		return nil, err
	}
	var entries []JournalEntry
	for _, e := range parsed {
		entries = append(entries, JournalEntry{
			Timestamp:    e.Timestamp,
			PriorityName: e.PriorityName,
			Message:      e.Message,
			Unit:         e.Unit,
			Hostname:     e.Hostname,
		})
	}
	return entries, nil
}

func (m *JournalModel) FetchWithFilter(count int, priority, unit string) ([]JournalEntry, error) {
	args := []string{"--no-pager", "--output=json", "-n"}
	c := strconv.Itoa(int(count))
	args = append(args, c)
	if priority != "" {
		args = append(args, "-p", priority)
	}
	if unit != "" && unit != "all" {
		args = append(args, "-u", unit)
	}
	out, err := exec.Command("journalctl", args...).Output()
	if err != nil {
		return nil, err
	}
	parsed, err := journalctl.Parse(strings.NewReader(string(out)))
	if err != nil {
		return nil, err
	}
	var entries []JournalEntry
	for _, e := range parsed {
		entries = append(entries, JournalEntry{
			Timestamp:    e.Timestamp,
			PriorityName: e.PriorityName,
			Message:      e.Message,
			Unit:         e.Unit,
			Hostname:     e.Hostname,
		})
	}
	return entries, nil
}
