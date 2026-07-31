# How to Build a New View

This guide documents the patterns and conventions for adding a new view to SystemPanel, based on the full chat history and all changes made.

## Architecture Overview

Every view follows this layered structure:

```
view/<name>/<name>.go    ←  GTK UI, signal handlers, widget tree
model/<name>.go           ←  system interaction, data fetching
parsers/<name>/<name>.go  ←  CLI output parsing into Go structs
detect/detect.go          ←  runtime dependency checks
main.go                   ←  registration in the sidebar
```

Optional layers that may be needed:
```
bindings/gtk4/            ←  new GTK4 CGo wrappers for missing widgets
widget/<name>.go          ←  reusable dialogs (sudo, connection editor)
```

## Step-by-Step: Adding a New View

### 1. Runtime Detection (`detect/detect.go`)

Each view registers a `DetectFn` that checks for required system tools and hardware. Add helper functions to `detect/detect.go`:

```go
// Check for a binary in PATH
func HasProgram(name string) bool {
    _, err := exec.LookPath(name)
    return err == nil
}

// Check for hardware via /sys
func HasWiFiHardware() bool {
    entries, _ := filepath.Glob("/sys/class/net/wl*")
    return len(entries) > 0
}
```

Add any new tool detection to `RunAll()` results map.

Detection runs once at startup. Views whose `DetectFn` returns `false` appear grayed out in the sidebar with "N/A" status and are unclickable. Views returning `true` load normally.

### 2. Parser (`parsers/<tool>/<tool>.go`)

Create a package that parses the CLI tool's output into public Go structs. The parser should be pure — no side effects, no `os/exec`, just string → struct conversion.

```go
package xrandr

import (
    "bufio"
    "strings"
)

type Output struct {
    Name      string
    Connected bool
    Modes     []Mode
}

type Mode struct {
    Width    int
    Height   int
    Refresh  float64
    Current  bool
}

func Parse(output string) (*Screen, []Output, error) {
    // Parse the CLI output line by line
    scanner := bufio.NewScanner(strings.NewReader(output))
    for scanner.Scan() {
        // ...
    }
    return &screen, outputs, scanner.Err()
}
```

**Convention**: One `Parse()` function that takes a `string` or `io.Reader` and returns parsed structs + error.

### 3. Model (`model/<name>.go`)

The model runs CLI commands and uses the parser. It provides data to the view.

```go
package model

type Monitor struct {
    Name       string
    Connected  bool
    Resolution string
    Modes      []ResolutionMode
}

type MonitorModel struct {
    observers []Observer
}

func (m *MonitorModel) ListMonitors() ([]Monitor, error) {
    out, err := exec.Command("xrandr").Output()
    if err != nil {
        return nil, err
    }
    _, outputs, err := xrandr.Parse(string(out))
    // Convert parser types to model types
    var monitors []Monitor
    for _, o := range outputs {
        monitors = append(monitors, Monitor{Name: o.Name, ...})
    }
    return monitors, nil
}

func (m *MonitorModel) SetResolution(output, mode string) error {
    return exec.Command("xrandr", "--output", output, "--mode", mode).Run()
}
```

**Convention**: Model types are simple structs with only the fields the view needs. Methods that modify system state are named as verbs: `SetX`, `ToggleX`, `Connect`, `Forget`.

### 4. View (`view/<name>/<name>.go`)

The view is the most complex layer. Follow this structure:

#### 4a. Descriptor

```go
var Descriptor = view.ViewDescriptor{
    Name:     "monitors",
    Title:    "Monitors",
    IconName: "video-display-symbolic",
    DetectFn: func() bool {
        return detect.HasProgram("xrandr")
    },
    Factory: func() view.View { return NewMonitorsView() },
}
```

- `Name`: unique ID, used for stack switching and settings toggles
- `Title`: displayed in the sidebar
- `IconName`: GTK icon name (use `-symbolic` variants)
- `DetectFn`: runtime check
- `Factory`: creates the view

#### 4b. Struct

```go
type MonitorsView struct {
    box         *gtk4.Box
    model       *model.MonitorModel
    listBox     *gtk4.ListBox
    connectBtn  *gtk4.Button      // optional: action button
    spinner     *gtk4.Spinner     // optional: loading indicator
    items       []monitorItem     // data items
    rowPtrToIdx map[unsafe.Pointer]int  // for row selection lookups
    rows        []*gtk4.ListBoxRow      // for cleanup on refresh
    parentWin   *gtk4.Window            // for dialogs
}
```

#### 4c. Constructor `NewXxxView()`

```go
func NewMonitorsView() *MonitorsView {
    mv := &MonitorsView{
        box:         gtk4.BoxNew(gtk4.OrientationVertical, 12),
        model:       &model.MonitorModel{},
        rowPtrToIdx: make(map[unsafe.Pointer]int),
    }
    mv.box.SetMarginStart(24)
    mv.box.SetMarginEnd(24)
    mv.box.SetMarginTop(24)
    mv.box.SetMarginBottom(24)

    // Header
    header := gtk4.LabelNew("Title")
    header.AddCSSClass("header-label")
    mv.box.Append(&header.Widget)

    // ListBox
    mv.listBox = gtk4.ListBoxNew()
    mv.listBox.SetSelectionMode(gtk4.SelectionSingle)
    // If selection is needed:
    mv.listBox.OnRowActivated(mv.onRowSelected)

    // ScrolledWindow
    scrollW := gtk4.ScrolledWindowNew()
    scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
    scrollW.SetHExpand(true)
    scrollW.SetVExpand(true)
    scrollW.SetChild(&mv.listBox.Widget)
    mv.box.Append(&scrollW.Widget)

    // Button bar at bottom
    btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
    btnBox.SetMarginTop(4)
    refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
    refreshBtn.OnClicked(func() { mv.refresh() })
    btnBox.Append(&refreshBtn.Widget)
    btnBox.Append(&mv.spinner.Widget)  // optional
    spacer := gtk4.LabelNew("")
    spacer.SetHExpand(true)
    btnBox.Append(&spacer.Widget)
    actionBtn := gtk4.ButtonNewWithLabel("Action")
    actionBtn.SetSensitive(false)
    actionBtn.OnClicked(mv.onAction)
    btnBox.Append(&actionBtn.Widget)
    mv.box.Append(&btnBox.Widget)

    mv.refresh()  // or mv.loadInitial() for async
    return mv
}
```

**Key layout pattern**: Header → ScrolledWindow(ListBox) → Button bar.

#### 4d. Refresh / Load

For fast synchronous loads, call directly:

```go
func (v *View) refresh() {
    // Clear old rows
    for _, r := range v.rows {
        v.listBox.Remove(r)
    }
    v.rows = v.rows[:0]
    v.items = v.items[:0]
    v.rowPtrToIdx = make(map[unsafe.Pointer]int)

    // Fetch data
    items, _ := v.model.ListItems()

    // Add section headers if needed
    if hasCategory1 {
        hdr := v.createSectionLabel("─ Category 1 ─")
        v.listBox.Append(hdr)
        v.rows = append(v.rows, hdr)
    }
    // Add data rows
    for _, item := range items {
        row := v.createRow(item)
        v.listBox.Append(row)
        v.rows = append(v.rows, row)
        v.rowPtrToIdx[row.Widget.Ptr()] = len(v.items)
    }
}
```

For slow operations (Bluetooth scan, network scan), use async refresh:

```go
func (v *View) refresh() {
    v.connectBtn.SetSensitive(false)
    gtk4.IdleAdd(func() { v.spinner.Start() })

    go func() {
        items, _ := v.model.FetchSlow()

        gtk4.IdleAdd(func() {
            v.spinner.Stop()
            v.populateList(items)
        })
    }()
}
```

**Important**: Never update GTK widgets from a goroutine — use `gtk4.IdleAdd(fn)` to schedule work on the GTK main thread.

#### 4e. Section Headers

```go
func (v *View) createSectionLabel(text string) *gtk4.ListBoxRow {
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

#### 4f. Data Rows

```go
func (v *View) createRow(item *itemType) *gtk4.ListBoxRow {
    row := gtk4.ListBoxRowNew()
    row.AddCSSClass("device-row")

    hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

    // Icon
    icon := gtk4.ImageNewFromIconName("icon-name-symbolic")
    icon.SetPixelSize(20)
    hbox.Append(&icon.Widget)

    // Main label (expands)
    nameLabel := gtk4.LabelNew(item.Name)
    nameLabel.SetHExpand(true)
    nameLabel.SetHAlign(gtk4.AlignStart)
    hbox.Append(&nameLabel.Widget)

    // Optional: status, percentage, buttons
    // ...

    row.SetChild(&hbox.Widget)
    return row
}
```

#### 4g. Row Selection (for action button pattern)

When you have an action button that enables on selection:

```go
func (v *View) onRowSelected(row *gtk4.ListBoxRow) {
    v.selectedID = ""
    v.connectBtn.SetSensitive(false)

    idx, ok := v.rowPtrToIdx[row.Widget.Ptr()]
    if ok && idx < len(v.items) {
        v.selectedID = v.items[idx].id
        v.connectBtn.SetSensitive(true)
    }
}

func (v *View) onActionClicked() {
    if v.selectedID == "" {
        return
    }
    v.model.DoSomething(v.selectedID)
    v.refresh()
}
```

**Critical**: `row.Widget.Ptr()` comparison uses the underlying C pointer. The signal handler creates a fresh Go `*ListBoxRow` wrapper, so direct `==` comparison on Go pointers will fail. Always use `rowPtrToIdx` or a similar map keyed by `Widget.Ptr()`.

#### 4h. Interface Methods

```go
func (v *View) Widget() *gtk4.Widget { return &v.box.Widget }
func (v *View) Name() string          { return "monitors" }
func (v *View) Title() string         { return "Monitors" }
func (v *View) IconName() string      { return "video-display-symbolic" }
func (v *View) OnShow()               { v.refresh() }
func (v *View) OnHide()               {}
func (v *View) Destroy()              {}
```

### 5. Registration (`main.go`)

Add the import and descriptor to the registry:

```go
import "github.com/cookiengineer/systempanel/view/monitors"

func registerViews() {
    view.Registry = []view.ViewDescriptor{
        // ...
        monitors.Descriptor,
        // ...
    }
}
```

### 6. New GTK Widgets (`bindings/gtk4/`)

If you need a GTK widget not yet wrapped:

1. Add the C wrapper in `bridge.c`:
   ```c
   void *gtk4MyWidgetNew(void) { return my_widget_new(); }
   void gtk4MyWidgetSetProperty(void *w, int val) { my_widget_set_property((MyWidget*)w, val); }
   ```

2. Declare the extern in `gtk4.go` preamble:
   ```go
   extern void *gtk4MyWidgetNew(void);
   extern void gtk4MyWidgetSetProperty(void *w, int val);
   ```

3. Add the Go type in `gtk4.go`:
   ```go
   type MyWidget struct{ Widget }
   func MyWidgetNew() *MyWidget { return &MyWidget{Widget{ptr: C.gtk4MyWidgetNew()}} }
   func (w *MyWidget) SetProperty(val int) { C.gtk4MyWidgetSetProperty(w.ptr, C.int(val)) }
   ```

## Common Patterns

### Form Dialogs (Connection Settings)

Used for Wi-Fi and LAN connection editing. Pattern:
- `GtkWindow` with `set_modal(true)` and `set_transient_for(parent)`
- Top bar with checkboxes and device selectors
- `GtkStackSwitcher` for tab switching
- `GtkStack` with one page per tab
- Bottom bar with Remove/Save buttons

### Sudo / Privileged Operations

See [SUDO-DIALOG.md](SUDO-DIALOG.md) for the complete integration guide, including all patterns, the three-tier escalation chain, and critical rules for goroutine safety.

Quick reference:

```go
// Pattern 2: Button-triggered privileged command (most common for views)
widget.PromptForSudo(view.parentWin, "Setting system time requires root privileges.", func(password string) {
    if password == "" {
        return
    }
    go func() {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        c := exec.CommandContext(ctx, "sudo", "-S", "-k", "timedatectl", "set-time", datetime)
        c.Stdin = strings.NewReader(password + "\n")
        c.Run()
        gtk4.IdleAdd(func() { view.loadAsync() })
    }()
})

// Model-layer: transparent escalation with cached password or pkexec fallback
widget.RunSudoCommand("cp", tmpPath, systemPath)
```

### Reusable Dialogs

Put reusable dialogs in `widget/`. Use the callback pattern to avoid blocking the GTK main loop:

```go
type MyDialog struct {
    win      *gtk4.Window
    callback func(result string)
}

func (d *MyDialog) Show(callback func(result string)) {
    d.callback = callback
    d.win.Present()
}
```

## View Checklist

When adding a new view, verify:

- [ ] `parsers/<tool>/` package with `Parse()` function
- [ ] `detect/detect.go` has tool detection helper
- [ ] `model/<name>.go` with Model struct and methods
- [ ] `view/<name>/<name>.go` with Descriptor and View implementation
- [ ] Section headers use `createSectionLabel("─ Name ─")` pattern
- [ ] Row selection uses `rowPtrToIdx` map with `Widget.Ptr()` keys
- [ ] Long operations use goroutine + `gtk4.IdleAdd()`
- [ ] Spinner shows during async loads
- [ ] GTK widgets never created/modified from goroutines
- [ ] `main.go` imports and registers the descriptor
- [ ] Clean build with `go build -o systempanel .`
