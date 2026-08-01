# GTK4 Style Guide

This document defines the UI and UX conventions for SystemPanel views. When creating a new view, follow these patterns to maintain consistency across the application.

## View Structure

Every view implements the `view.View` interface and registers via `view.ViewDescriptor`.

### File Layout

```
view/<name>/<name>.go    ←  GTK widget tree, signal handlers, action callbacks
model/<name>.go           ←  system interaction, data fetching
parsers/<name>/<name>.go  ←  CLI output parsing into Go structs
detect/detect.go          ←  runtime dependency checks
main.go                   ←  registration in sidebar
```

### Descriptor

```go
var Descriptor = view.ViewDescriptor{
    Name:     "mysection",
    Title:    "My Section",
    IconName: "my-icon-symbolic",
    DetectFn: func() bool {
        return detect.HasProgram("required-tool")
    },
    Factory: func() view.View { return NewMyView() },
}
```

- `Name`: unique ID, used for stack switching and settings toggles
- `Title`: displayed in the sidebar
- `IconName`: GTK icon name, always use `-symbolic` variants
- `DetectFn`: runtime presence check — returns `false` grays out the sidebar entry
- `Factory`: creates the view

### Struct

```go
type MyView struct {
    box         *gtk4.Box
    model       *model.MyModel
    listBox     *gtk4.ListBox
    rows        []*gtk4.ListBoxRow
    parentWin   *gtk4.Window
}
```

Required fields:
- `box` — the root widget, a vertical `Box` with orientation `gtk4.OrientationVertical`
- `parentWin` — set via `SetParentWindow()`, used by dialogs

Common optional fields:
- `listBox` — for list-based views
- `rows` — tracks appended rows for cleanup during refresh
- `items` — data items keyed by row, for row selection lookups
- `rowPtrToIdx` — `map[unsafe.Pointer]int` for `OnRowActivated` selection
- `spinner` — loading indicator
- `model` — data model reference

### Constructor

```go
func NewMyView() *MyView {
    mv := &MyView{
        box:   gtk4.BoxNew(gtk4.OrientationVertical, 12),
        model: &model.MyModel{},
    }
    mv.box.SetMarginStart(24)
    mv.box.SetMarginEnd(24)
    mv.box.SetMarginTop(24)
    mv.box.SetMarginBottom(24)

    // 1. Header
    header := gtk4.LabelNew("My Section")
    header.AddCSSClass("header-label")
    mv.box.Append(&header.Widget)

    // 2. Description
    desc := gtk4.LabelNew("Brief explanation of this view's purpose.")
    desc.SetWrap(true)
    desc.SetMarginBottom(12)
    mv.box.Append(&desc.Widget)

    // 3. ListBox inside ScrolledWindow
    mv.listBox = gtk4.ListBoxNew()
    mv.listBox.SetSelectionMode(gtk4.SelectionNone) // or SelectionSingle

    scrollW := gtk4.ScrolledWindowNew()
    scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
    scrollW.SetHExpand(true)
    scrollW.SetVExpand(true)
    scrollW.SetChild(&mv.listBox.Widget)
    mv.box.Append(&scrollW.Widget)

    // 4. Button bar at bottom
    btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
    btnBox.SetMarginTop(4)

    refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
    refreshBtn.OnClicked(func() { mv.refresh() })
    btnBox.Append(&refreshBtn.Widget)

    mv.box.Append(&btnBox.Widget)

    // 5. Initial load
    mv.refresh()

    return mv
}
```

The vertical structure is always: **Header → Description → ScrolledWindow(ListBox) → ButtonBar**.

### Interface Methods

```go
func (mv *MyView) Widget() *gtk4.Widget { return &mv.box.Widget }
func (mv *MyView) Name() string          { return "mysection" }
func (mv *MyView) Title() string         { return "My Section" }
func (mv *MyView) IconName() string      { return "my-icon-symbolic" }
func (mv *MyView) OnShow()               { mv.refresh() }
func (mv *MyView) OnHide()               {}
func (mv *MyView) Destroy()              {}
func (mv *MyView) SetParentWindow(p *gtk4.Window) { mv.parentWin = p }
```

---

## Refresh Pattern

### Synchronous (fast data fetch)

```go
func (mv *MyView) refresh() {
    for _, r := range mv.rows {
        mv.listBox.Remove(r)
    }
    mv.rows = mv.rows[:0]

    data, err := mv.model.ListItems()
    if err != nil || len(data) == 0 {
        label := gtk4.LabelNew("No items found")
        label.SetSensitive(false)
        label.SetHAlign(gtk4.AlignCenter)
        row := gtk4.ListBoxRowNew()
        row.SetChild(&label.Widget)
        mv.listBox.Append(row)
        mv.rows = append(mv.rows, row)
        return
    }

    for _, item := range data {
        row := mv.createRow(item)
        mv.listBox.Append(row)
        mv.rows = append(mv.rows, row)
    }
}
```

### Asynchronous (slow: network scan, service listing)

```go
func (mv *MyView) refresh() {
    gtk4.IdleAdd(func() { mv.spinner.Start() })

    go func() {
        data, _ := mv.model.FetchSlow()

        gtk4.IdleAdd(func() {
            mv.spinner.Stop()
            mv.populateList(data)
        })
    }()
}
```

Key rule: `refresh()` is always synchronous from the caller's perspective. The async work lives in a goroutine, and the GTK main thread is re-entered via `gtk4.IdleAdd()`.

---

## Empty State

Every view must handle the case where no data is available. Use a centered, insensitive label:

```go
label := gtk4.LabelNew("No storage devices detected")
label.SetSensitive(false)
label.SetHAlign(gtk4.AlignCenter)
row := gtk4.ListBoxRowNew()
row.SetChild(&label.Widget)
mv.listBox.Append(row)
mv.rows = append(mv.rows, row)
```

---

## Section Headers

Use `createSectionLabel()` to group related items:

```go
func (mv *MyView) createSectionLabel(text string) *gtk4.ListBoxRow {
    row := gtk4.ListBoxRowNew()
    row.SetSensitive(false)
    label := gtk4.LabelNew(text)
    label.SetHAlign(gtk4.AlignCenter)
    label.SetMarginTop(8)
    label.SetMarginBottom(4)
    label.SetSensitive(false)
    row.SetChild(&label.Widget)
    return row
}
```

Section header text uses em-dash separators: `"─ Connected Devices ─"`.

Always store section rows separately from data rows:

```go
mv.sectionRows = append(mv.sectionRows, hdr)
// vs
mv.rows = append(mv.rows, row)
```

Section rows and data rows must all be tracked for cleanup during refresh.

---

## Row Creation

### Simple Row (icon + label + optional widget)

Used by: batteries, bluetooth, icons, services, autostart.

```go
func (mv *MyView) createRow(item *Item) *gtk4.ListBoxRow {
    row := gtk4.ListBoxRowNew()
    row.AddCSSClass("device-row")

    hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

    icon := gtk4.ImageNewFromIconName("icon-name-symbolic")
    icon.SetPixelSize(20)
    hbox.Append(&icon.Widget)

    nameLabel := gtk4.LabelNew(item.Name)
    nameLabel.SetHExpand(true)
    nameLabel.SetHAlign(gtk4.AlignStart)
    hbox.Append(&nameLabel.Widget)

    if item.Status != "" {
        statusLabel := gtk4.LabelNew(item.Status)
        statusLabel.SetSensitive(false)
        hbox.Append(&statusLabel.Widget)
    }

    row.SetChild(&hbox.Widget)
    return row
}
```

### Grid Row (multiple aligned columns)

Used by: disks, wifi.

```go
grid := gtk4.GridNew()
grid.SetColumnSpacing(8)
grid.SetRowSpacing(2)

// Column 0: expanding name
nameLabel := gtk4.LabelNew(item.Name)
nameLabel.SetHExpand(true)
nameLabel.SetHAlign(gtk4.AlignStart)
grid.Attach(&nameLabel.Widget, 0, row, 1, 1)

// Column 1: fixed-width, right-aligned
sizeLabel := gtk4.LabelNew(item.Size)
sizeLabel.SetHAlign(gtk4.AlignEnd)
grid.Attach(&sizeLabel.Widget, 1, row, 1, 1)

// Column 2: fixed-width
typeLabel := gtk4.LabelNew(item.Type)
typeLabel.SetHAlign(gtk4.AlignStart)
grid.Attach(&typeLabel.Widget, 2, row, 1, 1)

// Column 3: action button
btn := gtk4.ButtonNewWithLabel("Action")
btn.SetSizeRequest(100, -1)
btn.SetHAlign(gtk4.AlignEnd)
grid.Attach(&btn.Widget, 3, row, 1, 1)
```

---

## Grid Alignment Rules

### Same-width widgets in a column

All widgets that share a Grid column across multiple rows **must** use the same `SetSizeRequest(width)`. Otherwise the Grid sizes each row independently and columns misalign.

```go
// Correct: both buttons share the same width
unmountBtn.SetSizeRequest(100, -1)
mountBtn.SetSizeRequest(100, -1)

// Wrong: different widths cause misalignment
unmountBtn.SetSizeRequest(80, -1)
mountBtn.SetSizeRequest(100, -1)
```

### Always attach a widget to every column

Even if a cell is empty, attach an empty placeholder so the column width is preserved:

```go
fsBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 0)
fsBox.SetSizeRequest(64, -1)

if fsText != "" {
    fsLabel := gtk4.LabelNew(fsText)
    fsBox.Append(&fsLabel.Widget)
}
// Always attach the box, even when empty
grid.Attach(&fsBox.Widget, col, row, 1, 1)
```

Without this, columns collapse when a row has no widget in that position, breaking alignment with other rows.

### Wrapping widgets in container boxes

If you wrap a button in a `Box` container, the container width comes from the button. Both must be consistent:

```go
btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 0)
btnBox.SetHAlign(gtk4.AlignEnd)
btnBox.Append(&button.Widget)
grid.Attach(&btnBox.Widget, col, row, 1, 1)
```

The `SetHAlign` on the box positions the button within the column; the box width is determined by the button's `SetSizeRequest`.

---

## Action Buttons

### In-row buttons (per-item actions)

Used by: disks (Mount/Unmount), bluetooth (Forget), services (gear).

```go
btn := gtk4.ButtonNewWithLabel("Action")
btn.SetSizeRequest(100, -1)
btn.SetHAlign(gtk4.AlignEnd)
btn.OnClicked(func() {
    id := item.ID
    // ... perform action ...
    mv.refresh()
})
```

### Bottom-bar buttons (global actions)

Used by: wifi (Connect), bluetooth (Connect), timedate (Save).

```go
actionBtn := gtk4.ButtonNewWithLabel("Connect")
actionBtn.SetSensitive(false)
actionBtn.OnClicked(mv.onActionClicked)
actionBtn.AddCSSClass("suggested-action") // optional, for primary actions
btnBox.Append(&actionBtn.Widget)
```

Bottom-bar layout: `[Refresh] [Spinner] [—spacer—] [Action]`

---

## Sudo Dialog Integration

See [SUDO-DIALOG.md](SUDO-DIALOG.md) for the complete guide.

### Pattern 2: Button-triggered privileged command

Used by: disks (unmount), wifi/lan/bluetooth (start daemon), timedate (set time).

```go
widget.PromptForSudo(mv.parentWin,
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
            gtk4.IdleAdd(func() { mv.loadAsync() })
        }()
    },
)
```

### NEVER capture parentWin at construction time

`parentWin` is nil during `NewXxxView()`. Always read it in the click callback closure:

```go
// Wrong: captured at construction, nil
pWin := mv.parentWin
btn.OnClicked(func() { widget.PromptForSudo(pWin, ...) })

// Correct: read at click time
btn.OnClicked(func() { widget.PromptForSudo(mv.parentWin, ...) })
```

---

## Row Selection

Used by views with `SelectionSingle` and a bottom-bar action button (wifi, bluetooth, services, icons).

```go
// Constructor
mv.listBox.OnRowActivated(mv.onRowSelected)
mv.rowPtrToIdx = make(map[unsafe.Pointer]int)

// During row creation
mv.rowPtrToIdx[row.Widget.Ptr()] = len(mv.items)

// Selection handler
func (mv *MyView) onRowSelected(row *gtk4.ListBoxRow) {
    mv.selectedID = ""
    mv.actionBtn.SetSensitive(false)

    idx, ok := mv.rowPtrToIdx[row.Widget.Ptr()]
    if ok && idx < len(mv.items) {
        mv.selectedID = mv.items[idx].ID
        mv.actionBtn.SetSensitive(true)
    }
}
```

**Critical**: Always key by `row.Widget.Ptr()` (the underlying C pointer), never by Go pointer identity. The signal handler receives a fresh Go `*ListBoxRow` wrapper, so `==` on Go pointers fails.

---

## Daemon Check Pattern

For views that depend on a system service (NetworkManager, bluetoothd):

```go
func (mv *MyView) checkDaemonAndConfigureButton() {
    if mv.model.IsServiceRunning() {
        mv.refreshBtn.SetLabel("Refresh")
        mv.refreshBtn.OnClicked(func() { mv.refresh() })
        mv.refresh()
    } else {
        mv.refreshBtn.SetLabel("Start Daemon")
        mv.refreshBtn.OnClicked(func() { mv.startDaemon() })
    }
}

func (mv *MyView) startDaemon() {
    widget.PromptForSudo(mv.parentWin, "Daemon requires root privileges to start.", func(password string) {
        if password == "" { return }
        go func() {
            exec.Command("sudo", "-S", "-k", "systemctl", "start", "myservice.service")
            gtk4.IdleAdd(func() { mv.checkDaemonAndConfigureButton() })
        }()
    })
}
```

The Refresh button is repurposed as "Start Daemon" when the service is not running, then reverts after startup.

---

## CSS Classes

### Always use `header-label` for view titles

```go
header := gtk4.LabelNew("Title")
header.AddCSSClass("header-label")
```

### Always use `device-row` for list rows

```go
row.AddCSSClass("device-row")
```

### Available commonly used CSS classes from `css/css.go`

| Class | Purpose |
|---|---|
| `header-label` | View title (18px, bold) |
| `device-row` | ListBox row padding |
| `monospace-label` | Monospace font for aligned data |
| `suggested-action` | Primary action button styling |
| `settings-row` | Settings toggle row padding |
| `group-header` | 15px bold with margin-top |
| `tag-pill` | Rounded tag with subtle background |

---

## Font and Typography

Use CSS classes, not inline markup, for font styling:

```go
label.AddCSSClass("monospace-label")
```

The `monospace-label` class is defined in `css/css.go` as `font-family: monospace;`.

---

## GTK Thread Safety

- **Never** create or modify GTK widgets from a goroutine
- **Never** call `time.Sleep` on the GTK main thread
- Always use `gtk4.IdleAdd(fn)` to schedule UI updates from goroutines
- Collect form values on the main thread before launching a goroutine

```go
btn.OnClicked(func() {
    // Read form values on the GTK main thread
    value := entry.GetText()

    go func() {
        // Slow work in goroutine
        result := doWork(value)

        gtk4.IdleAdd(func() {
            // Update UI on the main thread
            label.SetText(result)
        })
    }()
})
```

---

## Button Bar Layout

Every view has a bottom button bar with this pattern:

```
[Refresh]  [Spinner]  [——————————————— spacer ———————————————]  [Action]
```

The spacer is always:

```go
spacer := gtk4.LabelNew("")
spacer.SetHExpand(true)
btnBox.Append(&spacer.Widget)
```

---

## Spin Button Initial Values

When using `SpinButton`, set the display text before connecting `OnValueChanged` to avoid the initial value change firing:

```go
spin := gtk4.SpinButtonNew(1, 12, 1)
spin.SetValue(float64(now.Month()))
spin.SetText(fmt.Sprintf("%02.0f", float64(now.Month())))
spin.OnValueChanged(func(v float64) {
    spin.SetText(fmt.Sprintf("%02.0f", v))
})
```

---

## Icon Naming

Always use `-symbolic` icon variants for consistency with GTK4 symbolic icon themes:

- `drive-harddisk-symbolic`
- `network-wireless-symbolic`
- `bluetooth-active-symbolic`
- `battery-good-symbolic`
- `display-brightness-symbolic`
- `emblem-system-symbolic`
- `application-x-theme-symbolic`
- `preferences-system-time-symbolic`
- `applications-system-symbolic`
