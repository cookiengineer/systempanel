# System Panel

An adaptive GTK4 system control panel built with Go and CGo. Detects available system tools at runtime and presents only the views you can actually use.

![Demo](./docs/demo.gif)

## Features

- **Adaptive sidebar** — views for missing dependencies appear grayed out rather than hidden
- **16 built-in views**: LAN, Wi-Fi, Bluetooth, Volume, Monitors, Wallpapers, Brightness, Services, Autostart, Journal, Battery, Power Profile, Time & Date, Themes, Icons, Settings
- **Daemon auto-start** — Wi-Fi, LAN, Bluetooth, and Volume views detect inactive services and offer one-click start buttons
- **Wi-Fi connection editor** — create/edit NetworkManager profiles with tabbed settings (Wi-Fi, Security, IPv4, IPv6)
- **LAN connection management** — list and connect to Ethernet profiles
- **Backlight control** — slider-based brightness adjustment per display via `brightnessctl`
- **Power profile switching** — switch between power-saver, balanced, and performance modes via `powerprofilesctl`
- **Time & Date** — set system time, date, timezone, and NTP toggle via `timedatectl`
- **GTK theme switcher** — browse and apply GTK4/GTK3 themes instantly
- **Icon theme switcher** — browse and apply icon themes instantly
- **Monitor arrangement** — set resolution and relative positioning for multi-monitor setups via `xrandr`
- **Wallpaper manager** — browse and set wallpapers from ~/Pictures/Wallpapers using `feh` or `hsetroot`
- **Service unit editor** — inspect and edit systemd unit files with sudo
- **Sudo dialog** — session-cached password with visibility toggle and pkexec fallback
- **Zero third-party Go dependencies** — pure stdlib + GTK4 via CGo

## Views

| View | Requires | Description |
|------|----------|-------------|
| LAN | `nmcli` | Ethernet connection management |
| Wi-Fi | `nmcli` + wireless hardware | Network scan, connect, profile editor |
| Bluetooth | `bluetoothctl` + hardware | Device scan, connect, disconnect, forget |
| Volume | `pactl` + pulseaudio.service | Input/output device sliders, mute toggles |
| Monitors | `xrandr` | Resolution selection, multi-monitor arrangement |
| Wallpapers | `feh` or `hsetroot` | Browse and set wallpapers from ~/Pictures/Wallpapers |
| Brightness | `brightnessctl` | Per-display backlight slider control |
| Services | `systemctl` + systemd | Unit list with filter, start/stop/edit, user/system mode |
| Autostart | `systemctl` + XDG autostart dir | Desktop file enable/disable switches |
| Journal | `journalctl` | Color-coded log viewer with priority and unit filter |
| Battery | `upower` + battery hardware | Percentage, charge state, remaining time |
| Power Profile | `powerprofilesctl` | Power-saver/balanced/performance mode switcher |
| Time & Date | `timedatectl` | System time, date, timezone, NTP toggle |
| Themes | `/usr/share/themes/` | GTK4/GTK3 theme browser with instant apply |
| Icons | `/usr/share/icons/` | Icon theme browser with instant apply |
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

At startup, System Panel runs detection for each view — checking for binaries in `PATH` (`exec.LookPath`), hardware via `/sys`, and service availability (e.g. `systemctl --user list-unit-files pulseaudio.service`). Views whose dependencies are met get a clickable sidebar row. Missing views get a grayed-out "N/A" row.

For Wi-Fi, LAN, Bluetooth, and Volume views, clicking the view checks whether the underlying system service is running (`NetworkManager.service`, `bluetooth.service`, `pulseaudio.service`). If not, the Refresh button is replaced with a "Start X" button that launches the service (with sudo or user session as appropriate).

The GUI is built on GTK4 via handwritten CGo bindings in `bindings/gtk4/`. Widgets such as `ComboBoxText`, `Entry` (with visibility toggle for password fields), `SpinButton`, `CheckButton`, `Switch`, `LevelBar`, and `StackSwitcher` are wrapped in Go types with C bridge functions. A custom `GtkListBox` serves as the sidebar, and a `GtkStack` switches content views. All system interaction goes through subprocess calls to standard Linux tools — no DBus libraries needed.

Settings are persisted as JSON in `$XDG_CONFIG_HOME/systempanel/settings.json`. View visibility can be toggled per view from the Settings panel. Sudo operations are handled via a one-time password prompt with session caching, falling back to `pkexec` when cached credentials are unavailable. The password entry includes a visibility toggle button.

## License

[GPL 3.0 or later](./GPL-3.0-or-later.txt)

