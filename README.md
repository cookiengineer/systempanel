# System Panel

An adaptive GTK4 system control panel built with Go and CGo. Detects available system tools at runtime and presents only the views you can actually use.

## Features

- **Adaptive sidebar** — views for missing dependencies appear grayed out rather than hidden
- **10 built-in views**: Power, Volume, Wi-Fi, Bluetooth, Battery, Display, Services, Journal, Autostart, LAN
- **Wi-Fi connection editor** — create/edit NetworkManager profiles with a tabbed settings dialog (Wi-Fi, Security, IPv4, IPv6)
- **LAN connection management** — list and connect to Ethernet profiles
- **Zero third-party Go dependencies** — pure stdlib + GTK4 via CGo

## Views

| View | Requires | Description |
|------|----------|-------------|
| Power | `systemctl` | Shutdown, Reboot, Suspend, Hibernate |
| Volume | `pactl` + PulseAudio | Output device sliders, mute toggles |
| Wi-Fi | `nmcli` + wireless hardware | Network scan, connect, profile editor |
| Bluetooth | `bluetoothctl` + hardware | Device scan, connect, disconnect |
| Battery | `upower` + battery hardware | Percentage, state, remaining time |
| Display | `brightnessctl` | Brightness sliders per display |
| Services | `systemctl` | systemd unit list, start/stop |
| Journal | `journalctl` | Color-coded log viewer |
| Autostart | `systemctl` + XDG autostart dir | Desktop file enable/disable |
| LAN | `nmcli` | Ethernet connection management |

## Build

```bash
# Requirements: Go 1.22+, GTK 4.0 development headers
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

At startup, System Panel runs detection for each view — checking for binaries in `PATH` (`exec.LookPath`) and hardware via `/sys`. Views whose dependencies are met get a clickable sidebar row. Missing views get a grayed-out "N/A" row.

The GUI is built on GTK4 via handwritten CGo bindings in `bindings/gtk4/`. A custom `GtkListBox` serves as the sidebar, and a `GtkStack` switches content views. All system interaction goes through subprocess calls to standard Linux tools — no DBus libraries needed.

## License

GPLv3
