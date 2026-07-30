package widget

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/parsers/desktop"
)

type DesktopDialog struct {
	win       *gtk4.Window
	stack     *gtk4.Stack
	entry     *desktop.DesktopEntry
	parentWin *gtk4.Window
	filePath  string
	isNew     bool
	entries   map[string]*gtk4.Entry
	checks    map[string]*gtk4.CheckButton
	combos    map[string]*gtk4.ComboBoxText
}

func NewDesktopDialog(parent *gtk4.Window, entry *desktop.DesktopEntry, filePath string) *DesktopDialog {
	d := &DesktopDialog{
		parentWin: parent,
		entry:     entry,
		filePath:  filePath,
		isNew:     filePath == "",
		entries:   make(map[string]*gtk4.Entry),
		checks:    make(map[string]*gtk4.CheckButton),
		combos:    make(map[string]*gtk4.ComboBoxText),
	}
	if d.entry == nil {
		d.entry = &desktop.DesktopEntry{
			Path:     "",
			Sections: make(map[string]desktop.DesktopSection),
		}
	}
	d.build()
	return d
}

func (d *DesktopDialog) build() {
	d.win = gtk4.WindowNew()
	title := "New Autostart Entry"
	if !d.isNew {
		title = "Edit Autostart Entry - " + filepath.Base(d.filePath)
	}
	d.win.SetTitle(title)
	d.win.SetDefaultSize(800, 600)
	d.win.SetModal(true)
	d.win.SetTransientFor(d.parentWin)

	vbox := gtk4.BoxNew(gtk4.OrientationVertical, 0)

	switcher := gtk4.StackSwitcherNew()
	switcher.SetMarginStart(12)
	switcher.SetMarginTop(4)
	vbox.Append(&switcher.Widget)

	d.stack = gtk4.StackNew()
	d.stack.SetTransitionType(gtk4.StackTransitionCrossfade)
	d.stack.SetTransitionDuration(120)
	d.stack.SetVHomogeneous(true)
	d.stack.SetHExpand(true)
	d.stack.SetVExpand(true)

	tab := d.buildSectionTab("Desktop Entry")
	d.stack.AddTitled(&tab.Widget, "main", "Desktop Entry")
	d.stack.SetVisibleChildName("main")
	switcher.SetStack(d.stack)

	vbox.Append(&d.stack.Widget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginStart(24)
	btnBox.SetMarginEnd(24)
	btnBox.SetMarginTop(12)
	btnBox.SetMarginBottom(12)

	cancelBtn := gtk4.ButtonNewWithLabel("Cancel")
	cancelBtn.OnClicked(func() { d.win.Close() })
	btnBox.Append(&cancelBtn.Widget)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	btnBox.Append(&spacer.Widget)

	saveBtn := gtk4.ButtonNewWithLabel("Save")
	saveBtn.AddCSSClass("suggested-action")
	saveBtn.OnClicked(d.onSave)
	btnBox.Append(&saveBtn.Widget)

	vbox.Append(&btnBox.Widget)

	d.win.SetChild(&vbox.Widget)
}

func (d *DesktopDialog) buildSectionTab(sectionName string) *gtk4.Box {
	box := gtk4.BoxNew(gtk4.OrientationVertical, 8)
	box.SetMarginStart(24)
	box.SetMarginEnd(24)
	box.SetMarginTop(12)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)

	innerBox := gtk4.BoxNew(gtk4.OrientationVertical, 4)

	keys := desktop.SortDesktopKeys(d.entry.SectionKeys(sectionName))
	if len(keys) == 0 {
		keys = desktop.DesktopKnownKeys()
	}

	keySet := make(map[string]bool)
	for _, k := range keys {
		if !keySet[k] {
			keySet[k] = true
		}
	}
	for _, k := range d.entry.SectionKeys(sectionName) {
		if !keySet[k] {
			keySet[k] = true
			keys = append(keys, k)
		}
	}
	keys = desktop.SortDesktopKeys(keys)

	for _, key := range keys {
		if !keySet[key] {
			continue
		}
		value := d.entry.SectionKey(sectionName, key)

		switch {
		case key == "Type":
			row := d.createTypeRow(key, value)
			innerBox.Append(&row.Widget)
		case key == "Categories":
			row := d.createCategoriesRow(key, value)
			innerBox.Append(&row.Widget)
		case isDesktopBool(key):
			row := d.createCheckRow(key, value)
			innerBox.Append(&row.Widget)
		default:
			row := d.createEntryRow(key, value)
			innerBox.Append(&row.Widget)
		}
	}

	scrollW.SetChild(&innerBox.Widget)
	box.Append(&scrollW.Widget)

	sw := scrollW.Widget
	sw.SetVExpand(true)

	return box
}

func (d *DesktopDialog) createEntryRow(key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := gtk4.LabelNew(key + ":")
	lbl.SetSizeRequest(200, -1)
	lbl.SetHAlign(gtk4.AlignEnd)
	lbl.SetMarginEnd(4)
	row.Append(&lbl.Widget)

	entry := gtk4.EntryNew()
	entry.SetText(value)
	entry.SetHExpand(true)
	row.Append(&entry.Widget)

	d.entries[key] = entry

	return row
}

func (d *DesktopDialog) createCheckRow(key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := gtk4.LabelNew(key + ":")
	lbl.SetSizeRequest(200, -1)
	lbl.SetHAlign(gtk4.AlignEnd)
	lbl.SetMarginEnd(4)
	row.Append(&lbl.Widget)

	check := gtk4.CheckButtonNew()
	check.SetActive(value == "true")
	row.Append(&check.Widget)

	d.checks[key] = check

	return row
}

func (d *DesktopDialog) createTypeRow(key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := gtk4.LabelNew(key + ":")
	lbl.SetSizeRequest(200, -1)
	lbl.SetHAlign(gtk4.AlignEnd)
	lbl.SetMarginEnd(4)
	row.Append(&lbl.Widget)

	combo := gtk4.ComboBoxTextNew()
	combo.Append("Application", "Application")
	combo.Append("Link", "Link")
	combo.Append("Directory", "Directory")
	combo.SetHExpand(true)

	if value != "" {
		for i := 0; i < 3; i++ {
			combo.SetActive(i)
			if combo.GetActiveID() == value {
				break
			}
		}
	}
	row.Append(&combo.Widget)

	d.combos[key] = combo

	return row
}

func (d *DesktopDialog) createCategoriesRow(key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := gtk4.LabelNew(key + ":")
	lbl.SetSizeRequest(200, -1)
	lbl.SetHAlign(gtk4.AlignEnd)
	lbl.SetMarginEnd(4)
	row.Append(&lbl.Widget)

	combo := gtk4.ComboBoxTextNewWithEntry()
	for _, cat := range xdgCategories() {
		combo.Append(cat, cat)
	}
	combo.SetHExpand(true)

	if value != "" {
		for i, cat := range xdgCategories() {
			if cat == value {
				combo.SetActive(i)
				break
			}
		}
		if combo.GetActive() < 0 {
			combo.Append(value, value)
			combo.SetActive(len(xdgCategories()))
		}
	}
	row.Append(&combo.Widget)

	d.combos[key] = combo

	return row
}

func (d *DesktopDialog) onSave() {
	for key, entry := range d.entries {
		value := strings.TrimSpace(entry.GetText())
		d.entry.Set(key, value)
	}
	for key, check := range d.checks {
		if check.GetActive() {
			d.entry.Set(key, "true")
		} else {
			d.entry.Set(key, "")
		}
	}
	for key, combo := range d.combos {
		id := combo.GetActiveID()
		if id == "" {
			id = combo.GetActiveText()
		}
		d.entry.Set(key, strings.TrimSpace(id))
	}

	content := d.entry.Serialize()

	path := d.filePath
	if path == "" {
		if name := d.entry.Get("Name"); name != "" {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, ".config", "autostart", sanitizeDesktopFilename(name)+".desktop")
		}
	}
	if path == "" {
		return
	}

	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)
	os.WriteFile(path, []byte(content), 0644)

	d.win.Close()
}

func sanitizeDesktopFilename(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '/' || r == 0 {
			return -1
		}
		return r
	}, s)
}

func (d *DesktopDialog) Present() {
	d.win.Present()
}

func isDesktopBool(key string) bool {
	bools := map[string]bool{
		"NoDisplay":             true,
		"Hidden":                true,
		"Terminal":              true,
		"StartupNotify":         true,
		"DBusActivatable":       true,
		"PrefersNonDefaultGPU":  true,
		"SingleMainWindow":      true,
	}
	return bools[key]
}

func xdgCategories() []string {
	return []string{
		"Audio", "AudioVideo", "Video",
		"Development", "Education",
		"Game", "Graphics",
		"Network", "Office",
		"Science", "Settings",
		"System", "Utility",
		"Qt", "GTK", "KDE", "GNOME",
	}
}
