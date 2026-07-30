# SystemPanel Implementation

## Overview

SystemPanel is a GTK4 system control panel app written in Go with CGo bindings. It detects available system tools at runtime and adapts its interface accordingly — views for which dependencies are missing appear grayed out in the sidebar rather than being hidden.

## Technology Stack

- **Language**: Go 1.26+
- **GUI**: GTK 4.22+ via CGo (`#cgo pkg-config: gtk4`)
- **Dependencies**: Go stdlib only. No third-party Go libraries.
- **System tools**: `pactl`, `nmcli`, `bluetoothctl`, `upower`, `brightnessctl`, `systemctl`, `journalctl`, `xrandr` — all detected at runtime via `exec.LookPath`

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
│   ├── volume/                     # PulseAudio Input/Output devices, mute toggle, sliders
│   ├── wifi/                       # nmcli WiFi networks + connection editor
│   ├── bluetooth/                  # bluetoothctl async scan + signal strength
│   ├── battery/                    # upower battery info
│   ├── monitors/                   # xrandr resolution + arrangement control
│   ├── services/                   # systemd unit management
│   ├── journal/                    # journalctl log viewer
│   ├── autostart/                  # XDG autostart .desktop manager
│   ├── lan/                        # Ethernet connection manager
│   └── settings/                   # View visibility toggles + about
├── model/
│   ├── model.go                    # Model interface
│   ├── power.go                    # systemctl power commands
│   ├── volume.go                   # pactl volume control (sinks + sources)
│   ├── wifi.go                     # nmcli WiFi scan/connect/configured profiles
│   ├── bluetooth.go                # bluetoothctl async scan loop + signal
│   ├── battery.go                  # upower battery parsing
│   ├── monitor.go                  # xrandr monitor listing + resolution/position
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
│   ├── nmconnection/                # NetworkManager .nmconnection keyfile parser/serializer
│   └── xrandr/                      # xrandr output parsers (screens, outputs, modes)
├── widget/
│   ├── sudodialog.go               # Callback-based sudo password dialog (session-cached)
│   └── connectiondialog.go          # Connection Settings dialog with GtkStackSwitcher tabs
├── assets/
│   ├── styles.css                   # GTK4 CSS stylesheet
│   └── systempanel.desktop          # Application launcher
└── docs/
    ├── IMPLEMENTATION.md            # This file
    └── VIEW-GUIDE.md                # Step-by-step guide for adding new views
```

## CGo Binding Strategy

The CGo bindings are in `bindings/gtk4/`:

- **`bridge.c`** — A standalone C file containing all GTK4 API wrapper functions and signal bridge callbacks. Each GTK4 function is wrapped in a `gtk4*` function that takes `void*` pointers instead of typed pointers, avoiding type conflicts in Go. Includes bridge callbacks for 8 signal types, signal connection dispatcher, and `g_idle_add` wrapper for posting work from goroutines to the GTK main thread.
- **`gtk4.go`** — The single CGo entry point with `#cgo pkg-config: gtk4`, `//export` Go callback functions, all Go widget types, and helper functions for signal connection. Widget types include: Widget, Application, Window, HeaderBar, Box, Button, Label, Image, Stack, StackPage, ListBox, ListBoxRow, ScrolledWindow, Scale, Switch, Entry, CSS, Spinner, LevelBar, EventControllerKey, ComboBoxText (with entry variant), CheckButton, SpinButton, StackSwitcher.

All C function calls go through `bridge.c` wrappers. Signal callbacks use `runtime/cgo.Handle` for type-safe Go↔C bridging. The `connectGSignal` dispatcher function in `bridge.c` selects the correct C bridge callback based on signal type. The `gtk4.IdleAdd(fn)` function uses `g_idle_add` to safely call back into the GTK main thread from goroutines.

## Sidebar Design

A custom `GtkListBox` on the left side provides the sidebar. Each row represents a view. When the underlying dependency is detected, the row is clickable and switches the content area (`GtkStack`) to the corresponding view via `SetVisibleChildName`. When the dependency is missing, the row is rendered insensitive (grayed out with "N/A" status).

Row selection uses C pointer comparison (`Widget.Ptr()`) via a `rowToName map[unsafe.Pointer]string`, because the signal handler creates a fresh Go `*ListBoxRow` wrapper — direct Go pointer comparison always fails.

The Settings button in the headerbar toggles between the Settings view and the previously active view.

**Sidebar order** (static, defined in `main.go`):

1. LAN
2. Wi-Fi
3. Bluetooth
4. Volume
5. Monitors
6. Services
7. Autostart
8. Journal
9. Battery
10. Power
11. Settings

## Runtime Detection

Each view registers a `DetectFn`. Detection uses `exec.LookPath` for binaries, `/sys` filesystem checks for hardware, and process tests for server availability.

| View | Detection |
|------|-----------|
| Power | `systemctl` + `/run/systemd/system` |
| Volume | `pactl` + PulseAudio server reachable (`pactl info`) |
| Wi-Fi | `nmcli` + `/sys/class/net/wl*` exists |
| Bluetooth | `bluetoothctl` + `/sys/class/bluetooth` |
| Battery | `upower` + `/sys/class/power_supply/BAT*` |
| Monitors | `xrandr` |
| Services | `systemctl` + `/run/systemd/system` |
| Journal | `journalctl` |
| Autostart | `systemctl` + `$XDG_CONFIG_HOME/autostart/` |
| LAN | `nmcli` |
| Settings | Always true |

## Async Refresh Pattern

Views with slow operations (Bluetooth 30s scan, WiFi scan) use goroutines with `gtk4.IdleAdd` to avoid blocking the GTK main loop:

1. Show spinner via `gtk4.IdleAdd(func() { spinner.Start() })`
2. Launch goroutine that runs the slow model operation
3. When complete, post results back via `gtk4.IdleAdd(func() { populateList(results) })`
4. Stop spinner, clear scanning flag

**Critical**: Never create, modify, or read GTK widgets from goroutines. Always use `gtk4.IdleAdd(fn)` to schedule work on the GTK main thread.

## Sudo / Privileged Operations

The sudo dialog (`widget/sudodialog.go`) uses a callback-based pattern to avoid main loop blocking:

- **Startup**: `PromptForSudo()` shows the dialog after the window appears, caching the password for the session
- **Validation**: `sudo -S -k true` verifies the password before caching
- **Usage**: `RunSudoCommand(cmd, args...)` runs commands with `sudo -S -k` using the cached password
- **Fallback**: If no cached password, falls back to `pkexec`, then to re-prompting the dialog
- **Session cache**: `sessionPassword` package variable persists across operations; `InvalidateSudo()` clears it

Connection profile saves use a three-tier approach:
1. Direct `os.WriteFile` (if user has permissions)
2. `RunSudoCommand("cp", tmpPath, systemPath)` (cached sudo or pkexec)
3. `SudoDialog.RunCommand(...)` (re-prompts for password)

## Connection Settings Dialog

For WiFi and LAN views, the Connection Settings dialog uses `GtkStack` + `GtkStackSwitcher` for proper GTK4 tabbed navigation (replacing earlier manual button tab simulation):

- **Wi-Fi tab**: SSID, Mode (Client/Hotspot/Ad-hoc), Band (Auto/5GHz/2.4GHz), BSSID, Cloned MAC (editable combo: Preserve/Permanent/Random/Stable + custom), MTU
- **Wi-Fi Security tab**: Security type (None, WPA/WPA2/WPA3 Personal, WPA Enterprise, Dynamic WEP, WEP, WPA3 SAE), Password
- **IPv4 Settings tab**: Method (DHCP/Manual/Link-Local/Shared/Disabled), Address, Netmask, Gateway, DNS, never-default toggle (with left spacer for alignment)
- **IPv6 Settings tab**: Method (Auto/Manual/Link-Local/Shared/Ignore/Disabled), Address, Prefix, Gateway, DNS, never-default toggle (with left spacer for alignment)

Top bar: Connect automatically checkbox + Device selector (populated from `nmcli device status`, with `selectComboOrAppend` fallback for missing interfaces). Bottom bar: Remove button (existing profiles only) + Save button.

All form labels are unified to 130px width with right-alignment via shared `addField`, `addComboRow`, and `addSpinRow` helpers.

The dialog adapts to connection type: WiFi/802.11 shows all four tabs; 802.3-ethernet shows only IPv4/IPv6 tabs. Profiles are saved to `/etc/NetworkManager/system-connections/`. The `parsers/nmconnection` package handles `.nmconnection` keyfile format serialization with `NewWiFiConnection()` and `NewEthernetConnection()` factory functions.

Tmp files are written to `/tmp/` (user-writable) then copied via sudo. Whitespace in filenames is preserved (only `/` and null bytes are stripped by `sanitizeFilename`).

## Views

### Power
4 large colored buttons: Shutdown, Reboot, Suspend, Hibernate. Runs `systemctl` commands.

### Volume
ListBox-based view with two sections: "─ Input Devices ─" and "─ Output Devices ─". Each device row has:
- **Mute toggle button** — shows speaker icon when unmuted, muted icon when muted. Toggles via `pactl set-sink-mute`/`set-source-mute`.
- **Device name** — description from PulseAudio
- **Volume slider** — controls output volume via `pactl set-sink-volume`. Input device sliders skip value-changed on construction.

Uses `parsers/pactl` with a generic `parseDevices()` function matching both "Sink #" and "Source #" prefixes. Refresh button reloads the device list. The `HasPulseAudio()` detector verifies the server is reachable by running `pactl info`.

### Wi-Fi
Async scan via `nmcli device wifi list` with goroutine + `gtk4.IdleAdd`. Results sorted by configured/unconfigured:
- "─ Saved Networks ─" header (if any profiles exist) — configured rows show a settings gear icon to edit
- "─ Available Networks ─" header — unconfigured rows

Bottom bar: Refresh button + disabled Connect button. Connect button enables on row selection. Clicking Connect runs `nmcli` directly for configured networks or opens the Connection Settings dialog for unconfigured ones. A spinner shows during the async scan. Section header rows are tracked in `sectionRows` slice and properly removed on each refresh.

### Bluetooth
Async scan loop with continuous device polling:
1. `loadInitial()` — fast `bluetoothctl devices` listing on startup/view show (no scan)
2. Refresh button triggers `startScan()` — powers on adapter, enables discoverable, starts `bluetoothctl scan on`
3. Background goroutine polls `bluetoothctl devices` + `bluetoothctl info <MAC>` every 2 seconds via `ScanLoop`, pushing results through `gtk4.IdleAdd`
4. After 30 seconds, scan stops and final poll runs. Spinner shows during scan.
5. Re-entrancy prevented via `scanning` mutex flag

RSSI signal strength shown as percentage next to device name (mapped from dBm via `RSSIToPercent`). `parseInt` handles the hex format from `bluetoothctl info` output (`RSSI: 0xffffffa5 (-91)`).

Section headers: "─ Paired Devices ─" and "─ Available Devices ─". Paired devices show a "Forget" button. Row selection enables Connect button (runs `bluetoothctl connect`).

### Monitors
Replaces the old Display view. Uses `xrandr` for detection and control. Single ListBox with two sections:

- **─ Resolution ─**: Per monitor — star icon for primary, port name (eDP1, HDMI1), current resolution, and a resolution dropdown populated from available modes with refresh rates. Preferred and current modes are indicated.
- **─ Arrangement ─**: Only shown when 2+ monitors are connected. Per monitor — port name, relation dropdown (Left of/Right of/Above/Below/Same as), and target monitor dropdown (other connected monitors). When only one monitor is connected, shows "Connect a second monitor to configure arrangement" message. Dropdowns are 100px wide, centered in the row.

Apply button runs both resolution changes (`xrandr --output <name> --mode <res>`) and position changes (`xrandr --output <name> --<relation> <target>`) in sequence.

### Battery
Displays battery percentage bar, charge state, remaining time, and capacity. Uses `upower -i` for data.

### Services
Lists systemd units with active/failed/inactive status indicators. Filter entry for searching. Clicking a row toggles start/stop via `systemctl`.

### Journal
Fetches recent journal entries via `journalctl --output=json`. Color-coded by priority (emerg/err red, warning yellow, info green). Refresh and Clear buttons.

### Autostart
Lists `.desktop` files from XDG autostart directories with name, icon, executable. Toggle switch per entry to enable/disable via the `Hidden` key.

### LAN
Lists ethernet NetworkManager connections with settings gear icon per row. Bottom bar: Refresh + disabled Connect button. Connect button enables on row selection (runs `nmcli connection up`). When no connections exist, Connect button opens the Connection Settings dialog to create a new LAN profile (connType: `802-3-ethernet`).

### Settings
Toggle switches to show/hide individual views. About section with version info.
