# System Panel

An adaptive GTK4 system control panel built with Go and CGo. Detects available system tools at runtime and presents only the views you can actually use.

![Demo](./docs/demo.gif)

## Features

- **Adaptive sidebar** — views for missing dependencies appear grayed out rather than hidden
- **11 built-in views**: LAN, Wi-Fi, Bluetooth, Volume, Monitors, Services, Autostart, Journal, Battery, Themes, Settings
- **Wi-Fi connection editor** — create/edit NetworkManager profiles with tabbed settings (Wi-Fi, Security, IPv4, IPv6)
- **LAN connection management** — list and connect to Ethernet profiles
- **GTK theme switcher** — browse and apply GTK4/GTK3 themes instantly
- **Monitor arrangement** — set resolution and relative positioning for multi-monitor setups via `xrandr`
- **Service unit editor** — inspect and edit systemd unit files with sudo
- **Zero third-party Go dependencies** — pure stdlib + GTK4 via CGo

## Views

| View | Requires | Description |
|------|----------|-------------|
| LAN | `nmcli` | Ethernet connection management |
| Wi-Fi | `nmcli` + wireless hardware | Network scan, connect, profile editor |
| Bluetooth | `bluetoothctl` + hardware | Device scan, connect, disconnect, forget |
| Volume | `pactl` + PulseAudio server | Input/output device sliders, mute toggles |
| Monitors | `xrandr` | Resolution selection, multi-monitor arrangement |
| Services | `systemctl` + systemd | Unit list with filter, start/stop/edit, user/system mode |
| Autostart | `systemctl` + XDG autostart dir | Desktop file enable/disable switches |
| Journal | `journalctl` | Color-coded log viewer with priority and unit filter |
| Battery | `upower` + battery hardware | Percentage, charge state, remaining time |
| Themes | `/usr/share/themes/` | GTK4/GTK3 theme browser with instant apply |
| Settings | Always available | View visibility toggles, dark mode, about |

## Build

```bash
# Requirements: Go 1.26+, GTK 4.22+ development headers
make build
```

## Install

```bash
sudo make install
```

## Run

```bash
make run
# or
./systempanel
```

## How It Works

At startup, System Panel runs detection for each view — checking for binaries in `PATH` (`exec.LookPath`), hardware via `/sys`, and server availability (e.g. `pactl info`). Views whose dependencies are met get a clickable sidebar row. Missing views get a grayed-out "N/A" row.

The GUI is built on GTK4 via handwritten CGo bindings in `bindings/gtk4/`. Widgets such as `ComboBoxText`, `Entry` (with completion), `SpinButton`, `CheckButton`, `Switch`, `LevelBar`, and `StackSwitcher` are wrapped in Go types with C bridge functions. A custom `GtkListBox` serves as the sidebar, and a `GtkStack` switches content views. All system interaction goes through subprocess calls to standard Linux tools — no DBus libraries needed.

Settings are persisted as JSON in `$XDG_CONFIG_HOME/systempanel/settings.json`. View visibility can be toggled per view from the Settings panel. Sudo operations are handled via a one-time password prompt with session caching, falling back to `pkexec` when cached credentials are unavailable.

## License

GPLv3
