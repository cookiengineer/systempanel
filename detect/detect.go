package detect

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// HasProgram checks if a binary is available in PATH.
func HasProgram(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// HasWiFiHardware checks for wireless network interfaces.
func HasWiFiHardware() bool {
	entries, err := filepath.Glob("/sys/class/net/wl*")
	if err != nil {
		return false
	}
	return len(entries) > 0
}

// HasBatteryHardware checks for battery power supplies.
func HasBatteryHardware() bool {
	entries, err := filepath.Glob("/sys/class/power_supply/BAT*")
	if err != nil {
		return false
	}
	return len(entries) > 0
}

// HasBluetoothHardware checks for Bluetooth hardware.
func HasBluetoothHardware() bool {
	entries, err := filepath.Glob("/sys/class/bluetooth/hci*")
	if err != nil {
		return false
	}
	if len(entries) > 0 {
		return true
	}
	_, err = os.Stat("/sys/class/bluetooth")
	return err == nil
}

// HasAutostartDir checks if XDG autostart directories exist.
func HasAutostartDir() bool {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return false
		}
		configDir = filepath.Join(home, ".config")
	}
	autostartDir := filepath.Join(configDir, "autostart")
	info, err := os.Stat(autostartDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// HasDesktopSession checks if XDG_CURRENT_DESKTOP is set.
func HasDesktopSession() string {
	return strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))
}

// IsWayland checks if running under Wayland.
func IsWayland() bool {
	return os.Getenv("WAYLAND_DISPLAY") != ""
}

// HasSystemd checks if systemd is the init system.
func HasSystemd() bool {
	_, err := os.Stat("/run/systemd/system")
	return err == nil
}

// HasPulseAudio checks if the PulseAudio server is reachable via pactl.
func HasPulseAudio() bool {
	cmd := exec.Command("pactl", "info")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

// DetectionResult holds all runtime detection results.
type DetectionResult struct {
	Programs map[string]bool
	Hardware map[string]bool
	IsWayland bool
	Desktop   string
}

// RunAll performs all runtime detection checks.
func RunAll() DetectionResult {
	return DetectionResult{
		Programs: map[string]bool{
			"systemctl":       HasProgram("systemctl"),
			"journalctl":      HasProgram("journalctl"),
			"pactl":           HasProgram("pactl"),
			"nmcli":           HasProgram("nmcli"),
			"bluetoothctl":    HasProgram("bluetoothctl"),
			"upower":          HasProgram("upower"),
			"brightnessctl":   HasProgram("brightnessctl"),
			"powerprofilesctl": HasProgram("powerprofilesctl"),
			"gammastep":       HasProgram("gammastep"),
			"xrandr":          HasProgram("xrandr"),
		},
		Hardware: map[string]bool{
			"wifi":      HasWiFiHardware(),
			"battery":   HasBatteryHardware(),
			"bluetooth": HasBluetoothHardware(),
		},
		IsWayland: IsWayland(),
		Desktop:   HasDesktopSession(),
	}
}
