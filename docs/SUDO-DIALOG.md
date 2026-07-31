# Sudo Dialog Integration Guide

This document describes all patterns for integrating privileged operations using the Sudo Dialog system in `widget/sudodialog.go`.

## Core Concepts

The Sudo Dialog provides a GTK4 modal dialog that prompts the user for their password. Once authenticated, the password is cached in a session-scoped module variable (`sessionPassword`) so subsequent privileged operations do not re-prompt.

| Mechanism | Behavior |
|---|---|
| `PromptForSudo(parent, message, callback)` | Convenience: creates dialog, shows it, calls callback with password (or `""` if cancelled) |
| `sd.Show(callback)` | Shows the dialog unless password is already cached — then calls callback immediately |
| `sd.RunCommand(cmd, args...)` | Calls `Show()`, then runs the command via `sudo -S -k` in a goroutine |
| `RunSudoCommand(cmd, args...)` | Runs command via cached sudo password, or falls back to `pkexec` if no password cached |
| `InvalidateSudo()` | Clears the cached session password |
| `NewSudoDialog(parent)` | Low-level: creates dialog with default message |
| `NewSudoDialogWithMessage(parent, message)` | Low-level: creates dialog with custom message |

## Password Verification

Passwords are validated with `sudo -S -k true` (5s timeout). The `-k` flag ensures sudo's own credential cache does not interfere with the check. Invalid passwords show "Incorrect password. Please try again." and clear the entry field.

The dialog also supports Enter key to submit (keyval `0xff0d` via EventControllerKey) and a password visibility toggle button.

---

## Pattern 1: Session Pre-caching at App Startup

**File**: `app/panel.go:109`

Prompt for sudo immediately after the main window appears. This seeds the `sessionPassword` cache so all subsequent `RunSudoCommand()` calls and `PromptForSudo()` calls will skip the dialog (since the password is already cached).

```go
w := panel.win // *gtk4.Window

widget.PromptForSudo(w,
    "SystemPanel needs administrative privileges to manage network connections and system services.",
    func(password string) {
        if password != "" {
            log.Println("Sudo session cached for privileged operations")
        }
    },
)
```

**Key detail**: Even if the user cancels (password == `""`), the `PromptForSudo` / `RunSudoCommand` fallback chain still works correctly — `RunSudoCommand` falls back to `pkexec`, and `PromptForSudo` will show the dialog again on next use.

---

## Pattern 2: Button-Triggered Privileged Operation

Used by: Wi-Fi "Start NetworkManager", LAN "Start NetworkManager", Bluetooth "Start Bluetooth Daemon", Time/Date "Save".

This is the primary pattern for views. A button click triggers a privileged system command. The flow:

1. Collect data from form widgets on the GTK main thread
2. Call `widget.PromptForSudo()` with a descriptive message
3. In the callback, guard against cancellation (`password == ""`)
4. Run the command in a **goroutine** with `sudo -S -k`
5. Use `gtk4.IdleAdd()` to update the UI after the command completes

### Example: Wi-Fi Start NetworkManager (`view/wifi/wifi.go:131-147`)

```go
func (wv *WiFiView) startNetworkManager() {
    widget.PromptForSudo(wv.parentWin,
        "NetworkManager requires root privileges to start.",
        func(password string) {
            if password == "" {
                return
            }
            go func() {
                ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                defer cancel()
                c := exec.CommandContext(ctx, "sudo", "-S", "-k", "systemctl", "start", "NetworkManager.service")
                c.Stdin = strings.NewReader(password + "\n")
                c.Run()
                gtk4.IdleAdd(func() {
                    wv.checkDaemonAndConfigureButton()
                })
            }()
        },
    )
}
```

### Example: Time/Date Save (`view/timedate/timedate.go:88-123`)

```go
saveBtn := gtk4.ButtonNewWithLabel("Save")
saveBtn.OnClicked(func() {
    datetime := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
        int(tv.yearSpin.GetValue()),
        int(tv.monthSpin.GetValue()),
        int(tv.daySpin.GetValue()),
        int(tv.hourSpin.GetValue()),
        int(tv.minSpin.GetValue()),
        int(tv.secSpin.GetValue()),
    )
    zone := tv.tzCombo.GetActiveID()

    widget.PromptForSudo(tv.parentWin,
        "Setting system time requires root privileges.",
        func(password string) {
            if password == "" {
                return
            }
            go func() {
                ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                defer cancel()

                c := exec.CommandContext(ctx, "sudo", "-S", "-k", "timedatectl", "set-time", datetime)
                c.Stdin = strings.NewReader(password + "\n")
                c.Run()

                if zone != "" && zone != tv.model.CurrentTimezone() {
                    c := exec.CommandContext(ctx, "sudo", "-S", "-k", "timedatectl", "set-timezone", zone)
                    c.Stdin = strings.NewReader(password + "\n")
                    c.Run()
                }

                gtk4.IdleAdd(func() {
                    tv.loadAsync()
                })
            }()
        },
    )
})
```

### Pattern 2 Checklist

- [ ] Collect all form values **before** calling `PromptForSudo` (on the GTK main thread)
- [ ] Pass a descriptive message string as the second argument
- [ ] Guard against `password == ""` (user cancelled) in the callback
- [ ] Run `sudo -S -k <cmd>` in a **goroutine** (never block the GTK main thread)
- [ ] Use `context.WithTimeout` with 30s timeout
- [ ] Pipe the password to stdin via `strings.NewReader(password + "\n")`
- [ ] Use `gtk4.IdleAdd()` for any post-command UI updates
- [ ] Set the parent window for proper modal behavior

---

## Pattern 3: Transparent Escalation Chain (Three-Tier Fallback)

Used by: ConnectionDialog save (`widget/connectiondialog.go:519-541`), ServiceDialog save.

This pattern is for file write operations where the destination requires root. It silently tries three escalation levels before showing a dialog:

1. **Tier 1** — Direct write with `os.WriteFile`. Success means the user already has permission.
2. **Tier 2** — `widget.RunSudoCommand("cp", tmpPath, systemPath)`. Uses cached password via `sudo -S -k`, or falls back to `pkexec` (GUI password prompt).
3. **Tier 3** — `SudoDialog.RunCommand("cp", tmpPath, systemPath)`. Explicitly shows a fresh Sudo Dialog to re-prompt for password.

```go
func (d *ConnectionDialog) onSave() {
    d.writeToConnection()
    content := d.connection.Serialize()

    profileDir := "/etc/NetworkManager/system-connections"
    filename := sanitizeFilename(d.ssid) + ".nmconnection"
    if d.connPath != "" {
        filename = filepath.Base(d.connPath)
    }
    path := filepath.Join(profileDir, filename)

    if err := os.WriteFile(path, []byte(content), 0600); err != nil {
        tmpPath := filepath.Join("/tmp", filename+".tmp")
        os.WriteFile(tmpPath, []byte(content), 0600)
        if err := RunSudoCommand("cp", tmpPath, path); err != nil {     // Tier 2
            sd := NewSudoDialog(d.parentWin)                            // Tier 3
            sd.RunCommand("cp", tmpPath, path)
        }
        os.Remove(tmpPath)
    }
    exec.Command("nmcli", "connection", "reload").Run()
    d.win.Close()
}
```

Key details:
- The temporary file is always written to `/tmp/` (user-writable)
- `os.Remove(tmpPath)` in a defer-like position cleans up regardless of which tier succeeded
- `RunSudoCommand` returns an error — check it to decide whether to escalate to Tier 3
- `sd.RunCommand` runs asynchronously (command is spawned in a goroutine), so the dialog does not wait for the command to finish before closing

---

## Model-Layer Usage

Models sometimes call `widget.RunSudoCommand()` directly when the operation is triggered programmatically rather than through a user-facing button (e.g., NTP toggle via checkbox, timezone change via combo). This avoids the full `PromptForSudo` callback chain.

**File**: `model/timedate.go:45,58,62,66`

```go
func (m *TimeDateModel) SetTimezone(zone string) error {
    return widget.RunSudoCommand("timedatectl", "set-timezone", zone)
}

func (m *TimeDateModel) SetNTP(enabled bool) error {
    if enabled {
        return widget.RunSudoCommand("timedatectl", "set-ntp", "true")
    }
    return widget.RunSudoCommand("timedatectl", "set-ntp", "false")
}
```

## Direct pkexec Calls

Some destructive operations (e.g., removing connection profiles) bypass the Sudo Dialog entirely and use `pkexec` directly. This is appropriate for one-shot, non-repeating operations where the dialog overhead is undesirable.

**File**: `widget/connectiondialog.go:543-550`

```go
func (d *ConnectionDialog) onRemove() {
    if d.connPath == "" {
        return
    }
    exec.Command("pkexec", "rm", d.connPath).Run()
    exec.Command("nmcli", "connection", "reload").Run()
    d.win.Close()
}
```

---

## Summary: Which Pattern to Use

| Scenario | Pattern | API |
|---|---|---|
| App startup — cache password for session | Pattern 1 | `PromptForSudo(parent, msg, callback)` |
| View button triggers privileged command | Pattern 2 | `PromptForSudo(parent, msg, callback)` + goroutine + `sudo -S -k` |
| File save to system directory (e.g. `/etc/`) | Pattern 3 | Three-tier: `os.WriteFile` → `RunSudoCommand` → `NewSudoDialog.RunCommand` |
| Model method triggered by UI toggle (checkbox, combo) | Model-layer | `RunSudoCommand(cmd, args...)` |
| One-shot destructive operation | Direct pkexec | `exec.Command("pkexec", args...).Run()` |
| Force re-authentication | Invalidate cache | `InvalidateSudo()` |

## Critical Rules

1. **Never block the GTK main thread**. Always run `sudo` commands in a goroutine.
2. **Never access GTK widgets from goroutines**. Use `gtk4.IdleAdd()` for all UI updates after command completion.
3. **Always guard against empty password** in `PromptForSudo` callbacks (user hit Cancel).
4. **Use `parent` window** for modal dialog stacking. The parent window reference is stored on the view struct (e.g., `wv.parentWin`). This is set in the view constructor from the `Widget()` method or passed from the panel orchestrator.
5. **Temporary files go to `/tmp/`** for the three-tier escalation pattern — it is always user-writable.
6. **Use `context.WithTimeout`** with 30s for sudo commands to prevent hung goroutines.
