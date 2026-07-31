package powerprofile

import (
	"strings"
)

type Profile struct {
	Name        string
	Active      bool
	Driver      string
	Degraded    string
}

func ParseList(output string) []Profile {
	var profiles []Profile
	var current *Profile

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.HasSuffix(trimmed, ":") {
			if current != nil {
				profiles = append(profiles, *current)
			}
			name := strings.TrimSuffix(trimmed, ":")
			active := false
			if strings.HasPrefix(name, "*") {
				active = true
				name = strings.TrimPrefix(name, "*")
				name = strings.TrimSpace(name)
			}
			current = &Profile{
				Name:   name,
				Active: active,
			}
			continue
		}

		if current != nil {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(strings.ToLower(parts[0]))
				val := strings.TrimSpace(parts[1])
				switch key {
				case "driver":
					current.Driver = val
				case "degraded":
					current.Degraded = val
				}
			}
		}
	}

	if current != nil {
		profiles = append(profiles, *current)
	}

	return profiles
}
