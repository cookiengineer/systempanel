# SystemPanel Implementation

## Overview

SystemPanel is a GTK4 system control panel app written in Go with CGo bindings. It detects available system tools at runtime and adapts its interface accordingly — views for which dependencies are missing appear grayed out in the sidebar rather than being hidden.

## Technology Stack

- **Language**: Go 1.26+
- **GUI**: GTK 4.22+ via CGo (`#cgo pkg-config: gtk4`)
- **Dependencies**: Go stdlib only. No third-party Go libraries.
- **System tools**: `pactl`, `nmcli`, `bluetoothctl`, `upower`, `brightnessctl`, `systemctl`, `journalctl` — all detected at runtime via `exec.LookPath`

## Architecture

```
main.go → app/panel.go (SystemPanel orchestrator)
           ├── view/         (View interface + ViewDescriptors + widget trees)
           ├── model/        (Model interface + system state fetching)
           ├── controller/   (Controller interface + action handling)
           ├── detect/       (runtime dependency detection)
           ├── config/       (JSON settings persistence)
           ├── css/          (embedded CSS provider)
           ├── parsers/      (CLI output parsers with public structs)
           ├── widget/       (reusable dialogs: sudo prompt, connection settings)
           └── bindings/gtk4/ (CGo GTK4 widget wrappers)
```

## Project Structure

```
systempanel/
├── go.mod
├── main.go                         # Entry point + view registration
├── Makefile                        # build/install/run
├── bindings/
│   └── gtk4/
│       ├── bridge.c                # C wrapper functions (all GTK4 API calls)
│       └── gtk4.go                 # CGo preamble, Go exports, widget types
├── app/
│   └── panel.go                    # SystemPanel orchestrator (sidebar, stack, views)
├── view/
│   ├── view.go                     # View interface + ViewDescriptor registry
│   ├── power/                      # Shutdown/Reboot/Suspend/Hibernate
│   ├── volume/                     # PulseAudio sinks, sliders, mute
│   ├── wifi/                       # nmcli WiFi networks + connection editor
│   ├── bluetooth/                  # bluetoothctl device manager
│   ├── battery/                    # upower battery info
│   ├── display/                    # brightnessctl sliders
│   ├── services/                   # systemd unit management
│   ├── journal/                    # journalctl log viewer
│   ├── autostart/                  # XDG autostart .desktop manager
│   ├── lan/                        # Ethernet connection manager
│   └── settings/                   # View visibility toggles + about
├── model/
│   ├── model.go                    # Model interface
│   ├── power.go                    # systemctl power commands
│   ├── volume.go                   # pactl volume control
│   ├── wifi.go                     # nmcli WiFi scan/connect
│   ├── bluetooth.go                # bluetoothctl device control
│   ├── battery.go                  # upower battery parsing
│   ├── display.go                  # brightnessctl control
│   ├── services.go                 # systemctl unit listing
│   ├── journal.go                  # journalctl JSON log fetching
│   └── autostart.go                # Desktop file scanning
├── controller/
│   └── controller.go               # Controller interface
├── detect/
│   └── detect.go                   # exec.LookPath + /sys hardware checks
├── config/
│   └── config.go                   # Atomic JSON settings read/write
├── css/
│   └── css.go                      # Embedded GTK4 CSS provider
├── parsers/
│   ├── journalctl/                  # journalctl --output=json parsers
│   ├── pactl/                       # pactl list sinks/sources parsers
│   ├── nmcli/                       # nmcli WiFi scan parsers
│   ├── bluetoothctl/                # bluetoothctl info/devices parsers
│   ├── upower/                      # upower -i parsers
│   ├── brightnessctl/               # brightnessctl info parsers
│   ├── systemctl/                   # systemctl list-units parsers
│   ├── desktop/                     # XDG .desktop file parsers
│   └── nmconnection/                # NetworkManager .nmconnection keyfile parser/serializer
├── widget/
│   ├── sudodialog.go               # Reusable sudo password GTK dialog
│   └── connectiondialog.go          # Connection Settings dialog with tabs
└── assets/
    ├── styles.css                   # GTK4 CSS stylesheet
    └── systempanel.desktop          # Application launcher
```

## CGo Binding Strategy

The CGo bindings are in `bindings/gtk4/`:

- **`bridge.c`** — A standalone C file containing all GTK4 API wrapper functions and signal bridge callbacks. Each GTK4 function is wrapped in a `gtk4*` function that takes `void*` pointers instead of typed pointers, avoiding type conflicts in Go.
- **`gtk4.go`** — The single CGo entry point with `#cgo pkg-config: gtk4`, `//export` Go callback functions, all Go widget types (Widget, Window, Button, Label, Stack, ListBox, etc.), and helper functions for signal connection.

All C function calls go through `bridge.c` wrappers. Signal callbacks use `runtime/cgo.Handle` for type-safe Go↔C bridging.

## Sidebar Design

A custom `GtkListBox` on the left side provides the sidebar. Each row represents a view. When the underlying dependency is detected, the row is clickable and switches the content area (`GtkStack`) to the corresponding view. When the dependency is missing, the row is rendered insensitive (grayed out).

## Runtime Detection

Each view registers a `DetectFn`. Detection uses `exec.LookPath` for binaries and `/sys` filesystem checks for hardware.

| View | Detection |
|------|-----------|
| Power | `systemctl` + `/run/systemd/system` |
| Volume | `pactl` + PulseAudio server reachable |
| Wi-Fi | `nmcli` + `/sys/class/net/wl*` exists |
| Bluetooth | `bluetoothctl` + `/sys/class/bluetooth` |
| Battery | `upower` + `/sys/class/power_supply/BAT*` |
| Display | `brightnessctl` |
| Services | `systemctl` + `/run/systemd/system` |
| Journal | `journalctl` |
| Autostart | `systemctl` + `$XDG_CONFIG_HOME/autostart/` |
| LAN | `nmcli` |
| Settings | Always true |

## Connection Settings Dialog

For WiFi and LAN views, the Connection Settings dialog provides a tabbed interface (via `GtkStack` + `GtkStackSwitcher`) for editing NetworkManager connection profiles:

- **Wi-Fi tab**: SSID, Mode (Client/Hotspot/Ad-hoc), Band (Auto/5GHz/2.4GHz), BSSID, Cloned MAC (editable combo: Preserve/Permanent/Random/Stable + custom), MTU
- **Wi-Fi Security tab**: Security type (None, WPA/WPA2/WPA3 Personal, WPA Enterprise, Dynamic WEP, WEP, WPA3 SAE), Password
- **IPv4 Settings tab**: Method (DHCP/Manual/Link-Local/Shared/Disabled), Address, Netmask, Gateway, DNS, never-default toggle
- **IPv6 Settings tab**: Method (Auto/Manual/Link-Local/Shared/Ignore/Disabled), Address, Prefix, Gateway, DNS, never-default toggle

Top bar: Connect automatically checkbox + Device selector. Bottom bar: Remove button (for existing profiles) + Save button.

Profiles are saved to `/etc/NetworkManager/system-connections/` via `pkexec` or the reusable sudo dialog. The `parsers/nmconnection` package handles `.nmconnection` keyfile format serialization.

## Views

### Power
4 large colored buttons: Shutdown, Reboot, Suspend, Hibernate. Runs `systemctl` commands.

### Volume
Lists PulseAudio output sinks with mute button, device name, and volume slider. Calls `pactl set-sink-volume`/`set-sink-mute` on change.

### Wi-Fi
Scans networks via `nmcli device wifi list`. Configured profiles shown under "Saved Networks" header with a settings gear icon to edit. Unconfigured networks under "Available Networks". Bottom bar has Refresh button and a Connect button that enables on row selection. Clicking Connect runs `nmcli` directly for configured networks or opens the Connection Settings dialog for unconfigured ones.

### Bluetooth
Scans and lists Bluetooth devices. Shows paired/connected status. Buttons: Scan, Refresh. Clicking a row toggles connect/disconnect.

### Battery
Displays battery percentage bar, charge state, remaining time, and capacity. Uses `upower -i` for data.

### Display
Brightness slider per display device. Uses `brightnessctl` for get/set.

### Services
Lists systemd units with active/failed/inactive status indicators. Filter entry for searching. Clicking a row toggles start/stop via `systemctl`.

### Journal
Fetches recent journal entries via `journalctl --output=json`. Color-coded by priority (emerg/err red, warning yellow, info green). Refresh and Clear buttons.

### Autostart
Lists `.desktop` files from XDG autostart directories with name, icon, executable. Toggle switch per entry to enable/disable via the `Hidden` key.

### LAN
Lists ethernet NetworkManager connections with Connect button (enables on row selection). Settings gear icon per row opens the Connection Settings dialog. When no connections exist, Connect button opens the dialog to create a new LAN profile.

### Settings
Toggle switches to show/hide individual views. About section with version info.
