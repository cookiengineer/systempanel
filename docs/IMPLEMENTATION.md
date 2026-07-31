# SystemPanel Implementation

## Overview

SystemPanel is a GTK4 system control panel app written in Go with CGo bindings. It detects available system tools at runtime and adapts its interface accordingly — views for which dependencies are missing appear grayed out in the sidebar rather than being hidden.

## Technology Stack

- **Language**: Go 1.26+
- **GUI**: GTK 4.22+ via CGo (`#cgo pkg-config: gtk4`)
- **Dependencies**: Go stdlib only. No third-party Go libraries.
- **System tools**: `pactl`, `nmcli`, `bluetoothctl`, `upower`, `brightnessctl`, `systemctl`, `journalctl`, `xrandr`, `timedatectl`, `powerprofilesctl` — all detected at runtime via `exec.LookPath`

## Architecture

```
main.go → app/panel.go (SystemPanel orchestrator)
           ├── view/         (17 files, View interface + ViewDescriptors + widget trees)
           ├── model/        (15 files, Model interface + system state fetching)
           ├── controller/   (Controller interface, currently unused)
           ├── detect/       (runtime dependency detection, 15+ checks)
           ├── config/       (JSON settings persistence)
           ├── css/          (embedded CSS provider)
           ├── parsers/      (16 packages, CLI output parsers with public structs)
           ├── widget/       (4 reusable dialogs: sudo, connection, service, desktop)
           └── bindings/gtk4/ (2 files, CGo GTK4 widget wrappers)
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
│   ├── volume/                     # PulseAudio Input/Output devices, mute toggle, sliders
│   ├── wifi/                       # nmcli WiFi networks + connection editor
│   ├── bluetooth/                  # bluetoothctl async scan + signal strength
│   ├── batteries/                  # upower battery info
│   ├── disks/                       # lsblk + udisksctl + smartctl disk management
│   ├── monitors/                   # xrandr resolution + arrangement control
│   ├── brightness/                 # brightnessctl backlight slider control
│   ├── services/                   # systemd unit management
│   ├── journal/                    # journalctl log viewer
│   ├── autostart/                  # XDG autostart .desktop manager
│   ├── lan/                        # Ethernet connection manager
│   ├── powerprofile/               # powerprofilesctl profile switcher
│   ├── timedate/                   # system time, date, timezone, NTP
│   ├── themes/                     # GTK theme browser and switcher
│   ├── icons/                      # Icon theme browser and switcher
│   └── settings/                   # View visibility toggles + about
├── model/
│   ├── model.go                    # Model interface + Observer type
│   ├── volume.go                   # pactl volume control (sinks + sources)
│   ├── wifi.go                     # nmcli WiFi scan/connect/configured profiles
│   ├── bluetooth.go                # bluetoothctl async scan loop + signal
│   ├── battery.go                  # upower battery parsing
│   ├── monitor.go                  # xrandr monitor listing + resolution/position
│   ├── brightness.go               # brightnessctl backlight info + set
│   ├── services.go                 # systemctl unit listing + unit file I/O
│   ├── journal.go                  # journalctl JSON log fetching
│   ├── autostart.go                # Desktop file scanning + toggle
│   ├── powerprofile.go             # powerprofilesctl profile listing + switching
│   ├── timedate.go                 # timedatectl status, timezone, NTP, set-time
│   ├── theme.go                    # GTK theme discovery + apply (gsettings)
│   └── icon.go                     # Icon theme discovery + apply (gsettings)
├── controller/
│   └── controller.go               # Controller interface (unused, reserved for future)
├── detect/
│   └── detect.go                   # exec.LookPath + /sys hardware + desktop session checks
├── config/
│   └── config.go                   # Atomic JSON settings read/write
├── css/
│   └── css.go                      # Embedded GTK4 CSS provider (no external stylesheet)
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
│   ├── xrandr/                      # xrandr output parsers (screens, outputs, modes)
│   ├── gtktheme/                    # GTK theme index.theme parser
│   ├── powerprofile/                # powerprofilesctl list parsers
│   ├── systemdunit/                 # systemd unit file parser/serializer
│   └── timedatectl/                 # timedatectl status + list-timezones parsers
├── widget/
│   ├── sudodialog.go               # Callback-based sudo password dialog (session-cached)
│   ├── connectiondialog.go          # Connection Settings dialog with GtkStackSwitcher tabs
│   ├── servicedialog.go             # systemd unit file editor with advanced mode
│   └── desktopdialog.go            # .desktop file editor for autostart entries
├── assets/
│   ├── systempanel.desktop          # Application launcher (.desktop file)
│   └── (CSS styles are embedded in css/css.go)
└── docs/
    ├── IMPLEMENTATION.md            # This file
    ├── VIEW-GUIDE.md                # Step-by-step guide for adding new views
    └── SUDO-DIALOG.md               # Sudo Dialog integration patterns and API reference
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

**Sidebar order** (static, defined in `main.go:registerViews()`):

1. LAN
2. Wi-Fi
3. Bluetooth
4. Volume
5. Monitors
6. Wallpapers
7. Brightness
8. Batteries
9. Disks
10. Services
11. Autostart
12. Journal
13. Power Profile
14. Time & Date
15. Themes
16. Icons
17. Settings

## Runtime Detection

Each view registers a `DetectFn`. Detection uses `exec.LookPath` for binaries, `/sys` filesystem checks for hardware, and service availability tests.

| View | Detection |
|------|-----------|
| LAN | `nmcli` |
| Wi-Fi | `nmcli` + `/sys/class/net/wl*` exists |
| Bluetooth | `bluetoothctl` + `/sys/class/bluetooth` |
| Volume | `pactl` + PulseAudio server reachable (`pactl info`) |
| Monitors | `xrandr` |
| Wallpapers | `feh` or `hsetroot` |
| Brightness | `brightnessctl` |
| Services | `systemctl` + `/run/systemd/system` |
| Journal | `journalctl` |
| Autostart | `systemctl` + `$XDG_CONFIG_HOME/autostart/` |
| Battery | `upower` + `/sys/class/power_supply/BAT*` |
| Disks | `lsblk` + `udisksctl` |
| Power Profile | `powerprofilesctl` |
| Time & Date | `timedatectl` |
| Themes | `/usr/share/themes/*/gtk-4.0/gtk.css` or `gtk-3.0/gtk.css` |
| Icons | `/usr/share/icons/*/index.theme` |
| Settings | Always true |

## Async Refresh Pattern

Views with slow operations (Bluetooth 30s scan, WiFi scan) use goroutines with `gtk4.IdleAdd` to avoid blocking the GTK main loop:

1. Show spinner via `gtk4.IdleAdd(func() { spinner.Start() })`
2. Launch goroutine that runs the slow model operation
3. When complete, post results back via `gtk4.IdleAdd(func() { populateList(results) })`
4. Stop spinner, clear scanning flag

**Critical**: Never create, modify, or read GTK widgets from goroutines. Always use `gtk4.IdleAdd(fn)` to schedule work on the GTK main thread.

## Sudo / Privileged Operations

The sudo dialog (`widget/sudodialog.go`) uses a callback-based pattern to avoid main loop blocking. The full integration guide is in [SUDO-DIALOG.md](SUDO-DIALOG.md). Quick summary:

- **Startup**: `PromptForSudo()` shows the dialog after the window appears, caching the password for the session
- **Validation**: `sudo -S -k true` verifies the password before caching
- **Pattern 1 — Session pre-cache**: `PromptForSudo()` at app startup seeds the cache
- **Pattern 2 — Button-triggered**: `PromptForSudo(parent, msg, callback)` + goroutine + `sudo -S -k` (used by Wi-Fi Start, Time/Date Save, etc.)
- **Pattern 3 — Three-tier escalation**: `os.WriteFile` → `RunSudoCommand` (cached sudo or pkexec) → `SudoDialog.RunCommand` (re-prompts)
- **Model-layer**: `RunSudoCommand(cmd, args...)` for non-button operations (NTP toggle, timezone change)
- **Session cache**: `sessionPassword` package variable; `InvalidateSudo()` clears it

## Widget Dialogs

Beyond the Sudo Dialog, `widget/` contains three other reusable dialogs:

- **ConnectionDialog** (`widget/connectiondialog.go`) — Wi-Fi/LAN profile editor with `GtkStack` + `GtkStackSwitcher` for 4 tabbed pages (Wi-Fi, Wi-Fi Security, IPv4, IPv6). Adapts to connection type: WiFi shows all four tabs; Ethernet shows only IPv4/IPv6 tabs. Profiles save to `/etc/NetworkManager/system-connections/`.
- **ServiceDialog** (`widget/servicedialog.go`) — systemd unit file editor with basic/advanced mode toggle. Displays and edits the raw `.service` file content, with a toggle to expose all sections (Unit, Service, Install).
- **DesktopDialog** (`widget/desktopdialog.go`) — `.desktop` file editor for autostart entries. Edits Name, Comment, Exec, Icon fields.

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

### LAN
Lists ethernet NetworkManager connections with settings gear icon per row. Bottom bar: Refresh + disabled Connect button. Connect button enables on row selection (runs `nmcli connection up`). When no connections exist, Connect button opens the Connection Settings dialog to create a new LAN profile (connType: `802-3-ethernet`). When NetworkManager daemon is not running, the Refresh button is replaced with a "Start NetworkManager" button (uses Sudo Dialog Pattern 2).

### Wi-Fi
Async scan via `nmcli device wifi list` with goroutine + `gtk4.IdleAdd`. Results sorted by configured/unconfigured:

- "─ Saved Networks ─" header (if any profiles exist) — configured rows show a settings gear icon to edit
- "─ Available Networks ─" header — unconfigured rows

Bottom bar: Refresh button + disabled Connect button. Connect button enables on row selection. Clicking Connect runs `nmcli` directly for configured networks or opens the Connection Settings dialog for unconfigured ones. A spinner shows during the async scan. Section header rows are tracked in `sectionRows` slice and properly removed on each refresh.

When NetworkManager daemon is not running, the Refresh button is replaced with a "Start NetworkManager" button (uses Sudo Dialog Pattern 2).

### Bluetooth
Async scan loop with continuous device polling:
1. `loadInitial()` — fast `bluetoothctl devices` listing on startup/view show (no scan)
2. Refresh button triggers `startScan()` — powers on adapter, enables discoverable, starts `bluetoothctl scan on`
3. Background goroutine polls `bluetoothctl devices` + `bluetoothctl info <MAC>` every 2 seconds via `ScanLoop`, pushing results through `gtk4.IdleAdd`
4. After 30 seconds, scan stops and final poll runs. Spinner shows during scan.
5. Re-entrancy prevented via `scanning` mutex flag

RSSI signal strength shown as percentage next to device name (mapped from dBm via `RSSIToPercent`). `parseInt` handles the hex format from `bluetoothctl info` output (`RSSI: 0xffffffa5 (-91)`).

Section headers: "─ Paired Devices ─" and "─ Available Devices ─". Paired devices show a "Forget" button. Row selection enables Connect button (runs `bluetoothctl connect`).

### Volume
ListBox-based view with two sections: "─ Input Devices ─" and "─ Output Devices ─". Each device row has:
- **Mute toggle button** — shows speaker icon when unmuted, muted icon when muted. Toggles via `pactl set-sink-mute`/`set-source-mute`.
- **Device name** — description from PulseAudio
- **Volume slider** — controls output volume via `pactl set-sink-volume`. Input device sliders skip value-changed on construction.

Uses `parsers/pactl` with a generic `parseDevices()` function matching both "Sink #" and "Source #" prefixes. Refresh button reloads the device list. The `HasPulseAudio()` detector verifies the server is reachable by running `pactl info`. When the PulseAudio user service is not running, a "Start PulseAudio" button appears.

### Monitors
Uses `xrandr` for detection and control. Single ListBox with two sections:

- **─ Resolution ─**: Per monitor — star icon for primary, port name (eDP1, HDMI1), current resolution, and a resolution dropdown populated from available modes with refresh rates. Preferred and current modes are indicated.
- **─ Arrangement ─**: Only shown when 2+ monitors are connected. Per monitor — port name, relation dropdown (Left of/Right of/Above/Below/Same as), and target monitor dropdown (other connected monitors). When only one monitor is connected, shows "Connect a second monitor to configure arrangement" message. Dropdowns are 100px wide, centered in the row.

Apply button runs both resolution changes (`xrandr --output <name> --mode <res>`) and position changes (`xrandr --output <name> --<relation> <target>`) in sequence.

### Wallpapers
Sets the X11 root window wallpaper using `feh` or `hsetroot`. Scans `~/Pictures/Wallpapers/` for image files (jpg, jpeg, png, bmp, gif, webp). List items show a 64px thumbnail on the left and filename on the right. Bottom bar: Refresh button, Mode dropdown (scale/stretch/center/tile/max), and Save button. Shows placeholder messages when the directory is missing or contains no images.

### Brightness
Slider-based backlight control via `brightnessctl`. Lists each backlight device (e.g. `intel_backlight`, `acpi_video0`) with a scale widget showing percentage. Adjusting the slider calls `brightnessctl set <percentage>%` for the selected device. Refresh button reloads current values.

### Services
Lists systemd units with active/failed/inactive status indicators. Toggle between User and System unit mode. Filter entry for searching by name. Clicking a row opens the ServiceDialog editor for basic/advanced unit file inspection and modification. Start/Stop/Restart buttons for the selected unit. Reload systemd button after unit file saves.

### Autostart
Lists `.desktop` files from XDG autostart directories with name, icon, executable. Toggle switch per entry to enable/disable via the `Hidden` key. Clicking a row opens the DesktopDialog editor for modifying desktop entry fields.

### Journal
Fetches recent journal entries via `journalctl --output=json`. Color-coded by priority (emerg/err red, warning yellow, info green). Filter by priority and unit name. Refresh and Clear buttons. Inline log viewer with scrollable container.

### Batteries
Displays battery percentage via LevelBar, charge state, remaining time (to empty or full), capacity, model, and vendor. Uses `upower -i` for data. Supports multiple batteries.

### Disks
Lists connected storage devices (HDDs, SSDs, USB drives) with their partitions using `lsblk --json`. Shows device type, transport, serial, and human-readable sizes. Each partition row includes filesystem type, mount point, and action buttons:

- **Mount/Unmount** — uses `udisksctl mount -b` / `udisksctl unmount -b`
- **LUKS Unlock** — uses `udisksctl unlock -b` for encrypted volumes

The **Refresh** button rescans via `lsblk` (no sudo required). If `smartctl` is available, it then prompts the Sudo Dialog (Pattern 2) and fetches S.M.A.R.T. health data for all drives via `smartctl --xall --json`. Health status is color-coded:

- **Healthy** (green) — S.M.A.R.T. passed, no reallocated/pending/uncorrectable sectors
- **Warning** (yellow) — elevated temperature or non-zero critical attribute counts
- **FAILING** (red) — S.M.A.R.T. failed or attributes marked as `when_failed`

Health details include a LevelBar, power-on hours, temperature (with warning threshold from the drive), and critical S.M.A.R.T. attributes (Reallocated Sectors, Current Pending Sectors, Uncorrectable Sectors, Wear Leveling Count). Mount/unmount/unlock actions trigger a fast rescan preserving existing S.M.A.R.T. data.

### Power Profile
Radio-button-based profile switcher via `powerprofilesctl`. Shows available profiles (power-saver, balanced, performance) with the active profile indicated. Selecting a different profile calls `powerprofilesctl set <profile>`. Driver and degraded status shown when applicable.

### Time & Date
Form-based view with SpinButton fields for year, month, day, hour, minute, second. Current time displayed in a large header label. Timezone combo populated from `timedatectl list-timezones`. NTP toggle via CheckButton. Save button runs `timedatectl set-time` and optionally `timedatectl set-timezone` using the Sudo Dialog (Pattern 2).

### Themes
Lists GTK4 and GTK3 themes from `/usr/share/themes/` with current theme indicator. Clicking a theme applies it via `gsettings set org.gnome.desktop.interface gtk-theme`. Theme discovery uses the `parsers/gtktheme` package to parse `index.theme` files.

### Icons
Lists icon themes from `/usr/share/icons/` with current icon theme indicator. Clicking an icon theme applies it via `gsettings set org.gnome.desktop.interface icon-theme`. Theme discovery checks for `index.theme` files in icon directories.

### Settings
Toggle switches to show/hide individual views in the sidebar. About section with version info. Visibility changes are persisted to `$XDG_CONFIG_HOME/systempanel/settings.json` and applied immediately by refreshing sidebar sensitivity.
