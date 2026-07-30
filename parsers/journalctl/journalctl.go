package journalctl

import (
	"bufio"
	"encoding/json"
	"io"
	"time"
)

// JournalEntry represents a single journal log entry parsed from journalctl --output=json.
type JournalEntry struct {
	Timestamp       time.Time `json:"-"`
	RealTimeUsec    string    `json:"__REALTIME_TIMESTAMP"`
	MonotonicUsec   string    `json:"__MONOTONIC_TIMESTAMP"`
	BootID          string    `json:"_BOOT_ID"`
	UID             string    `json:"_UID"`
	GID             string    `json:"_GID"`
	MachineID       string    `json:"_MACHINE_ID"`
	Hostname        string    `json:"_HOSTNAME"`
	Transport       string    `json:"_TRANSPORT"`
	Priority        string    `json:"PRIORITY"`
	PriorityName    string    `json:"-"`
	Message         string    `json:"MESSAGE"`
	MessageID       string    `json:"MESSAGE_ID"`
	ErrNo           string    `json:"ERRNO"`
	SyslogFacility  string    `json:"SYSLOG_FACILITY"`
	SyslogIdentifier string   `json:"SYSLOG_IDENTIFIER"`
	Unit            string    `json:"_SYSTEMD_UNIT"`
	UserUnit        string    `json:"_SYSTEMD_USER_UNIT"`
	Slice           string    `json:"_SYSTEMD_SLICE"`
	Cmdline         string    `json:"_CMDLINE"`
	Exe             string    `json:"_EXE"`
	Comm            string    `json:"_COMM"`
	PID             string    `json:"_PID"`
	CodeFile        string    `json:"CODE_FILE"`
	CodeLine        string    `json:"CODE_LINE"`
	CodeFunc        string    `json:"CODE_FUNC"`
	Raw             map[string]json.RawMessage
}

// PriorityNames maps systemd priority numbers to human-readable labels.
var PriorityNames = map[string]string{
	"0": "emerg",
	"1": "alert",
	"2": "crit",
	"3": "err",
	"4": "warning",
	"5": "notice",
	"6": "info",
	"7": "debug",
}

// Parse reads journalctl JSON Lines output from r and returns parsed entries.
func Parse(r io.Reader) ([]JournalEntry, error) {
	var entries []JournalEntry
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<20)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var entry JournalEntry
		if err := json.Unmarshal(line, &entry.Raw); err != nil {
			return entries, err
		}
		entry.populate()
		entries = append(entries, entry)
	}
	return entries, scanner.Err()
}

func (e *JournalEntry) populate() {
	if e.Raw == nil {
		return
	}
	if v, ok := e.Raw["__REALTIME_TIMESTAMP"]; ok {
		var s string
		json.Unmarshal(v, &s)
		e.RealTimeUsec = s
		if ts, err := parseMicroseconds(s); err == nil {
			e.Timestamp = ts
		}
	}
	if v, ok := e.Raw["MESSAGE"]; ok {
		json.Unmarshal(v, &e.Message)
	}
	if v, ok := e.Raw["PRIORITY"]; ok {
		json.Unmarshal(v, &e.Priority)
		e.PriorityName = PriorityNames[e.Priority]
	}
	if v, ok := e.Raw["_HOSTNAME"]; ok {
		json.Unmarshal(v, &e.Hostname)
	}
	if v, ok := e.Raw["_SYSTEMD_UNIT"]; ok {
		json.Unmarshal(v, &e.Unit)
	}
	if v, ok := e.Raw["SYSLOG_IDENTIFIER"]; ok {
		json.Unmarshal(v, &e.SyslogIdentifier)
	}
	if v, ok := e.Raw["_PID"]; ok {
		json.Unmarshal(v, &e.PID)
	}
	if v, ok := e.Raw["_UID"]; ok {
		json.Unmarshal(v, &e.UID)
	}
	if v, ok := e.Raw["_GID"]; ok {
		json.Unmarshal(v, &e.GID)
	}
	if v, ok := e.Raw["_BOOT_ID"]; ok {
		json.Unmarshal(v, &e.BootID)
	}
	if v, ok := e.Raw["_CMDLINE"]; ok {
		json.Unmarshal(v, &e.Cmdline)
	}
	if v, ok := e.Raw["SYSLOG_FACILITY"]; ok {
		json.Unmarshal(v, &e.SyslogFacility)
	}
	e.Raw = nil
}

func parseMicroseconds(s string) (time.Time, error) {
	var usec int64
	for _, c := range s {
		if c < '0' || c > '9' {
			continue
		}
		usec = usec*10 + int64(c-'0')
	}
	return time.Unix(usec/1_000_000, (usec%1_000_000)*1000), nil
}
